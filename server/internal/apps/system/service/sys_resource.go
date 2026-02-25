package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/common"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/auth_engine"
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/consts"
	pkgdto "biz-auto-api/pkg/dto"
	pkgmodels "biz-auto-api/pkg/models"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

const (
	AuthAllResourceNo  int64 = 0
	AuthAllResourceYes int64 = 1
)

type SysResource struct {
	service.Service
}

func (s *SysResource) GetResourceList(req *dto.GetResourceListReq) (*dto.GetResourceListRes, error) {
	var (
		log       = s.Log
		db        = s.DB
		result    = &dto.GetResourceListRes{}
		total     int64
		resources = make([]*models.SysResourceType, 0)
		err       error
	)
	err = db.Model(models.SysResourceType{}).Preload("Fields").Scopes(func(tx *gorm.DB) *gorm.DB {
		if tools.ToValue(req.Name) != "" {
			tx.Where("name like ?", tools.FuzzyQuery(tools.ToValue(req.Name)))
		}
		if tools.ToValue(req.Identify) != "" {
			tx.Where("identify like ?", tools.FuzzyQuery(tools.ToValue(req.Identify)))
		}
		if tools.ToValue(req.Table) != "" {
			tx.Where("`table` like ?", tools.FuzzyQuery(tools.ToValue(req.Table)))
		}
		return tx
	}).Count(&total).Order("sort,id asc").Scopes(req.MakePagination()).Find(&resources).Error
	if err != nil {
		err = errors.Wrap(err, "list resources failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	items := make([]*dto.ResourceItem, 0, len(resources))
	err = copier.CopyWithOption(&items, resources, copier.Option{Converters: tools.WithConverts(tools.GetTime2StrPtrConvert())})
	if err != nil {
		err = errors.Wrap(err, "copy resources failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	result.Items = items
	result.Total = total

	return result, nil
}
func (s *SysResource) GetResourceTableList(req *dto.GetResourceTableListReq) (*dto.GetResourceTableFieldListRes, error) {
	var (
		db     = s.DB
		log    = s.Log
		err    error
		result = &dto.GetResourceTableFieldListRes{}
		conf   = config.SystemConfig
		total  int64
		tables = make([]string, 0)
	)
	err = db.Table("INFORMATION_SCHEMA.TABLES").Scopes(func(tx *gorm.DB) *gorm.DB {
		tx.Where("TABLE_SCHEMA = ?", conf.Mysql.Database)
		if tools.ToValue(req.Name) != "" {
			tx.Where("TABLE_NAME like ?", tools.FuzzyQuery(tools.ToValue(req.Name)))
		}
		return tx
	}).Count(&total).Order("TABLE_NAME asc").Scopes(req.MakePagination()).Pluck("TABLE_NAME", &tables).Error
	if err != nil {
		err = errors.Wrap(err, "list tables failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	for _, t := range tables {
		table := &dto.TableItem{
			Name: &t,
		}
		result.Items = append(result.Items, table)
	}
	result.Total = total
	return result, nil
}
func (s *SysResource) GetResourceTableField(req *dto.GetResourceTableFieldReq) (*dto.GetResourceTableFieldRes, error) {
	var (
		db     = s.DB
		log    = s.Log
		err    error
		result = &dto.GetResourceTableFieldRes{}
		conf   = config.SystemConfig
		fields = make([]string, 0)
	)
	err = db.Table("INFORMATION_SCHEMA.COLUMNS").Scopes(func(tx *gorm.DB) *gorm.DB {
		tx.Where("TABLE_SCHEMA = ?", conf.Mysql.Database)
		tx.Where("TABLE_NAME = ?", tools.ToValue(req.TableName))
		return tx
	}).Order("COLUMN_NAME ASC").Pluck("COLUMN_NAME", &fields).Error
	if err != nil {
		err = errors.Wrap(err, "query table field failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	result.Fields = fields
	return result, nil
}
func (s *SysResource) CheckTableAndField(tableName string, fields []string) error {
	// 校验查询表是否存在
	getTableReq := &dto.GetResourceTableListReq{
		Name: &tableName,
	}
	getTableRes, err := s.GetResourceTableList(getTableReq)
	if err != nil {
		err = errors.WithMessage(err, "get resource table failed")
		return err
	}
	exist := false
	for _, t := range getTableRes.Items {
		if tools.PtrEqVal(t.Name, tableName) {
			exist = true
			break
		}
	}
	if !exist {
		return errors.Errorf("table(%v) not found", tableName)
	}
	// 校验字段信息是否存在
	getFieldReq := &dto.GetResourceTableFieldReq{
		TableName: &tableName,
	}
	getFieldRes, err := s.GetResourceTableField(getFieldReq)
	if err != nil {
		return errors.WithMessage(err, "get resource table field failed")
	}
	for _, field := range fields {
		if !tools.InSlice(field, getFieldRes.Fields) {
			return errors.Errorf("field(%v) not exist", field)
		}
	}
	return nil
}
func (s *SysResource) CheckQueryData(tableName string, fields []string, filter string) error {
	var (
		db  = s.DB
		obj = make(map[string]interface{})
	)
	if len(fields) == 0 {
		return errors.New("fields can not be empty")
	}
	fs := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		fs = append(fs, f)
	}
	err := db.Table(tableName).Select(fs[0], fs[1:]...).Scopes(func(tx *gorm.DB) *gorm.DB {
		if len(filter) > 0 {
			return tx.Where(filter)
		}
		return tx
	}).Limit(1).Find(&obj).Error
	if err != nil {
		return errors.Wrap(err, "check query data failed")
	}
	return nil
}
func (s *SysResource) AddResource(req *dto.AddResourceReq) (*dto.AddResourceRes, error) {
	var (
		db              = s.DB
		log             = s.Log
		err             error
		result          = &dto.AddResourceRes{}
		userId          = s.GetCurrentUserId()
		permissionTypes = make([]*models.PermissionType, 0, len(req.PermissionTypes))
		fields          = make([]*models.SysResourceTypeField, 0, len(req.Fields))
	)

	err = s.CheckTableAndField(tools.ToValue(req.Table), tools.GetSlice(req.Fields, func(e *dto.SysResourceTypeField) string { return tools.ToValue(e.FieldName) }))
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	err = s.CheckQueryData(tools.ToValue(req.Table), tools.GetSlice(req.Fields, func(e *dto.SysResourceTypeField) string { return tools.ToValue(e.FieldName) }), tools.ToValue(req.Filter))
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	permissionTypeCodes := make(map[string]struct{})
	for _, permissionType := range req.PermissionTypes {
		if _, ok := permissionTypeCodes[tools.ToValue(permissionType.Code)]; ok {
			return nil, errors.Errorf("权限类型编码重复: %v", tools.ToValue(permissionType.Code))
		} else {
			permissionTypeCodes[tools.ToValue(permissionType.Code)] = struct{}{}
		}
	}

	// 权限类型
	for _, permissionType := range req.PermissionTypes {
		p := &models.PermissionType{
			Code: permissionType.Code,
			Name: permissionType.Name,
		}
		permissionTypes = append(permissionTypes, p)
	}
	// 资源字段
	for _, field := range req.Fields {
		f := &models.SysResourceTypeField{
			ColumnName:    field.ColumnName,
			DictKey:       field.DictKey,
			FieldName:     field.FieldName,
			IsDict:        field.IsDict,
			ShowWithTag:   field.ShowWithTag,
			SupportFilter: field.SupportFilter,
			Sort:          field.Sort,
		}
		fields = append(fields, f)
	}
	// 去除空格
	if len(tools.ToValue(req.Filter)) > 0 {
		req.Filter = tools.ToPointer(strings.TrimSpace(tools.ToValue(req.Filter)))
	}
	resource := &models.SysResourceType{
		Name:            req.Name,
		Identify:        req.Identify,
		Table:           req.Table,
		Filter:          req.Filter,
		PermissionTypes: permissionTypes,
		ControlBy: pkgmodels.ControlBy{
			CreateBy: userId,
			UpdateBy: userId,
		},
		Fields: fields,
		Sort:   req.Sort,
	}
	err = db.Create(resource).Error
	if err != nil {
		log.Errorf("%+v", err)
		return nil, errors.Wrap(err, "create resource failed")
	}
	return result, nil
}
func (s *SysResource) UpdateResource(req *dto.UpdateResourceReq) (*dto.UpdateResourceRes, error) {
	var (
		db     = s.DB
		log    = s.Log
		err    error
		result = &dto.UpdateResourceRes{}
		userId = s.GetCurrentUserId()
		fields = make([]*models.SysResourceTypeField, 0, len(req.Fields))
	)
	err = s.CheckTableAndField(tools.ToValue(req.Table), tools.GetSlice(req.Fields, func(e *dto.SysResourceTypeField) string { return tools.ToValue(e.FieldName) }))
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	err = s.CheckQueryData(tools.ToValue(req.Table), tools.GetSlice(req.Fields, func(e *dto.SysResourceTypeField) string { return tools.ToValue(e.FieldName) }), tools.ToValue(req.Filter))
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	permissionTypeCodes := make(map[string]struct{})
	for _, permissionType := range req.PermissionTypes {
		if _, ok := permissionTypeCodes[tools.ToValue(permissionType.Code)]; ok {
			return nil, errors.Errorf("权限类型编码重复: %v", tools.ToValue(permissionType.Code))
		} else {
			permissionTypeCodes[tools.ToValue(permissionType.Code)] = struct{}{}
		}
	}
	// 去除空格
	if len(tools.ToValue(req.Filter)) > 0 {
		req.Filter = tools.ToPointer(strings.TrimSpace(tools.ToValue(req.Filter)))
	}
	// 主表
	data := tools.StructToMap(req, "", true, "Id", "PermissionTypes", "Fields")
	data["update_by"] = userId
	bs, err := json.Marshal(req.PermissionTypes)
	if err != nil {
		err = errors.Wrap(err, "json marshal permission types failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	data["permission_types"] = string(bs)
	for _, field := range req.Fields {
		fields = append(fields, &models.SysResourceTypeField{
			ResourceTypeId: req.Id,
			ColumnName:     field.ColumnName,
			DictKey:        field.DictKey,
			FieldName:      field.FieldName,
			IsDict:         field.IsDict,
			ShowWithTag:    field.ShowWithTag,
			SupportFilter:  field.SupportFilter,
			Sort:           field.Sort,
		})
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		// 修改资源信息
		err = tx.Model(&models.SysResourceType{}).Where("id = ?", req.Id).Updates(data).Error
		if err != nil {
			return errors.Wrap(err, "update resource type failed")
		}
		// 删除之前的字段信息
		err = tx.Model(&models.SysResourceTypeField{}).Where("resource_type_id = ?", req.Id).Delete(&models.SysResourceTypeField{}).Error
		if err != nil {
			return errors.Wrap(err, "delete resource type fields failed")
		}
		// 插入新的字段信息
		err = tx.CreateInBatches(fields, 100).Error
		if err != nil {
			return errors.Wrap(err, "batch create resource type fields failed")
		}
		return nil
	})
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}
func (s *SysResource) GetResourceViewFormSchemas(req *dto.GetResourceViewFormSchemasReq) (*dto.GetResourceViewFormSchemasRes, error) {
	var (
		result   = &dto.GetResourceViewFormSchemasRes{}
		schemas  = make([]*dto.FormSchema, 0)
		db       = s.DB
		log      = s.Log
		resource = models.SysResourceType{}
	)
	err := db.Model(&models.SysResourceType{}).Preload("Fields", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort asc")
	}).Where("id = ?", req.Id).First(&resource).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	for _, field := range resource.Fields {
		schema, err := s.BuildFormSchema(field)
		if err != nil {
			err = errors.WithMessage(err, "build form schema failed")
			log.Errorf("%+v", err)
			return nil, err
		}
		if schema != nil {
			schemas = append(schemas, schema)
		}
	}
	result.Schemas = schemas
	return result, nil
}
func (s *SysResource) BuildFormSchema(field *models.SysResourceTypeField) (*dto.FormSchema, error) {
	if tools.PtrEqVal(field.SupportFilter, models.ResourceFieldSupportFilterNo) {
		return nil, nil
	}
	schema := &dto.FormSchema{
		Field: field.FieldName,
		Label: field.ColumnName,
		ColProps: &struct {
			Span *int64 `json:"span,omitempty"`
		}{
			Span: tools.ToPointer(int64(8)),
		},
		ComponentProps: map[string]interface{}{},
	}
	if tools.PtrEqVal(field.IsDict, models.ResourceFieldIsDictYes) {
		schema.Component = tools.ToPointer("Select")
		dicts, err := common.QueryDictByKeys(s.DB, []string{tools.ToValue(field.DictKey)})
		if err != nil {
			return nil, err
		}
		options := make([]interface{}, 0)
		if dictMap, exist := dicts[tools.ToValue(field.DictKey)]; exist {
			for dictValue, dictLabel := range dictMap {
				option := &struct {
					Label string `json:"label"`
					Value string `json:"value"`
				}{
					Label: dictLabel,
					Value: dictValue,
				}
				options = append(options, option)
			}
		}
		schema.ComponentProps["options"] = options
		schema.ComponentProps["mode"] = "multiple"
		schema.ComponentProps["maxTagCount"] = 2
	} else {
		schema.Component = tools.ToPointer("Input")
	}
	return schema, nil
}
func (s *SysResource) GetResourceViewTableColumns(req *dto.GetResourceViewTableColumnsReq) (*dto.GetResourceViewTableColumnsRes, error) {
	var (
		result         = &dto.GetResourceViewTableColumnsRes{}
		columns        = make([]*dto.TableColumn, 0)
		db             = s.DB
		log            = s.Log
		resource       = models.SysResourceType{}
		showEnumFields = make(map[string]string)
		showTagFields  = make(map[string]string)
	)
	err := db.Model(&models.SysResourceType{}).Preload("Fields", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort asc")
	}).Where("id = ?", req.Id).First(&resource).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	// 多添加一列id
	idColumn := &dto.TableColumn{
		Title:     tools.ToPointer("id"),
		Resizable: tools.ToPointer(true),
		DataIndex: tools.ToPointer("id"),
		Width:     tools.ToPointer(int64(80)),
	}
	columns = append(columns, idColumn)
	for _, field := range resource.Fields {
		column := s.BuildTableColumn(field)
		if column != nil {
			//// 最后一列自适应宽度
			//if i == len(resource.Fields)-1 {
			//	column.Resizable = tools.ToPointer(false)
			//	column.Width = nil
			//}
			columns = append(columns, column)
		}
		if !tools.PtrEqVal(field.IsDict, models.ResourceFieldIsDictYes) {
			continue
		}
		if tools.PtrEqVal(field.ShowWithTag, models.ResourceFieldShowWithTagYes) {
			showTagFields[tools.ToValue(field.FieldName)] = tools.ToValue(field.DictKey)
		} else {
			showEnumFields[tools.ToValue(field.FieldName)] = tools.ToValue(field.DictKey)
		}
	}

	result.Columns = columns
	result.ShowEnumFields = showEnumFields
	result.ShowTagFields = showTagFields
	return result, nil
}
func (s *SysResource) BuildTableColumn(field *models.SysResourceTypeField) *dto.TableColumn {
	return &dto.TableColumn{
		Title:     field.ColumnName,
		Resizable: tools.ToPointer(true),
		DataIndex: field.FieldName,
		Width:     tools.ToPointer(int64(100)),
	}
}
func (s *SysResource) GetResourceDetail(req *dto.GetResourceDetailListReq) (*dto.GetResourceDetailListRes, error) {
	var (
		result   = &dto.GetResourceDetailListRes{}
		db       = s.DB
		log      = s.Log
		resource = models.SysResourceType{}
		items    = make([]map[string]interface{}, 0)
		total    int64
	)
	req.ExtraParams = make(dto.GetResourceDetailListExtraParams)
	// 获取所有查询参数
	query := s.Ctx.Request.URL.Query()
	for key, values := range query {
		// 排除已知参数
		if key != "id" && key != "page" && key != "pageSize" && key != "pageIndex" {
			if len(values) == 1 {
				req.ExtraParams[key] = values[0] // 取第一个值
			} else if len(values) > 1 {
				req.ExtraParams[key] = values
			}
		}
	}
	err := db.Model(&models.SysResourceType{}).Preload("Fields").First(&resource, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	fieldMap := make(map[string]bool)
	fields := make([]interface{}, 0)
	fields = append(fields, "id")
	for _, field := range resource.Fields {
		fieldMap[tools.ToValue(field.FieldName)] = tools.PtrEqVal(field.IsDict, models.ResourceFieldIsDictYes)
		fields = append(fields, fmt.Sprintf("`%s`", tools.ToValue(field.FieldName)))
	}
	if len(fields) == 0 {
		return nil, errors.New("resource fields is empty")
	}
	err = db.Table(tools.ToValue(resource.Table)).
		Select(fields[0], fields[1:]...).Scopes(func(tx *gorm.DB) *gorm.DB {
		for paramKey, paramValue := range req.ExtraParams {
			paramKey = strings.TrimSuffix(paramKey, "[]")
			if isDict, exist := fieldMap[paramKey]; exist {
				// 字典类型使用全等
				if isDict {
					switch paramValue.(type) {
					case string:
						tx.Where(fmt.Sprintf("%s = ?", paramKey), paramValue)
					case []string:
						tx.Where(fmt.Sprintf("%s in ?", paramKey), paramValue)
					}
				} else {
					switch paramValue.(type) {
					case string:
						tx.Where(fmt.Sprintf("%s like ?", paramKey), tools.FuzzyQuery(paramValue.(string)))
					case []string:
						or := ""
						for i, _ := range paramValue.([]string) {
							if i > 0 {
								or += " or "
							}
							or += fmt.Sprintf("%s like ?", paramKey)
						}
						tx.Where(or, func() []string {
							paramValues := make([]string, 0, len(paramValue.([]string)))
							for _, v := range paramValue.([]string) {
								paramValues = append(paramValues, tools.FuzzyQuery(v))
							}
							return paramValues
						})
					}
				}
			}
		}
		if len(tools.ToValue(resource.Filter)) > 0 {
			tx.Where(tools.ToValue(resource.Filter))
		}
		return tx
	}).
		Count(&total).
		Order("id asc").
		Scopes(req.MakePagination()).
		Find(&items).Error
	if err != nil {
		err = errors.Wrap(err, "query resource table data failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	for _, item := range items {
		result.Items = append(result.Items, dto.ResourceDetailItem(item))
	}
	result.Total = total
	return result, nil
}
func (s *SysResource) DeleteResource(req *dto.DeleteResourceReq) (*dto.DeleteResourceRes, error) {
	var (
		result   = &dto.DeleteResourceRes{}
		db       = s.DB
		log      = s.Log
		resource = models.SysResourceType{}
	)
	err := db.First(&resource, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Unscoped().Delete(&models.SysResourceType{}, req.Id).Error
		if err != nil {
			return errors.Wrap(err, "delete resource failed")
		}
		err = tx.Unscoped().Model(&models.SysResourceTypeField{}).
			Where("resource_type_id = ?", req.Id).
			Delete(&models.SysResourceTypeField{}).Error
		if err != nil {
			return errors.Wrap(err, "delete resource fields failed")
		}
		err = tx.Unscoped().Model(&models.SysRoleResource{}).
			Where("resource_type_identify = ?", resource.Identify).
			Delete(&models.SysRoleResource{}).Error
		if err != nil {
			return errors.Wrap(err, "delete resource roles failed")
		}
		return nil
	})
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}
func (s *SysResource) GetRoleResourceInfo(req *dto.GetRoleResourceInfoReq) (*dto.GetRoleResourceInfoRes, error) {
	var (
		result                           = &dto.GetRoleResourceInfoRes{}
		db                               = s.DB
		log                              = s.Log
		resourceType                     = models.SysResourceType{}
		authedAllResourcePermissionTypes = make([]string, 0)
	)

	err := db.Scopes(func(tx *gorm.DB) *gorm.DB {
		if tools.ToValue(req.Identify) != "" {
			return tx.Where("identify = ?", req.Identify)
		}
		if tools.ToValue(req.ResourceTypeId) != 0 {
			return tx.Where("id = ?", req.ResourceTypeId)
		}
		return tx
	}).First(&resourceType).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	// 查询角色授权的信息
	err = db.Model(models.SysRoleResource{}).Distinct("action").Select("action").
		Where("resource_type_identify = ? and role_id = ? and resource_id = ? ", resourceType.Identify, req.RoleId, auth_engine.Wildcard).
		Pluck("action", &authedAllResourcePermissionTypes).Error
	if err != nil {
		err = errors.Wrap(err, "query resource role failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var permissionTypes = make([]*dto.PermissionType, 0, len(resourceType.PermissionTypes))
	for _, permissionType := range resourceType.PermissionTypes {
		permissionTypes = append(permissionTypes, &dto.PermissionType{
			Code: permissionType.Code,
			Name: permissionType.Name,
		})
	}
	result.PermissionTypes = permissionTypes
	result.AuthedAllResourcePermissionTypes = authedAllResourcePermissionTypes
	return result, nil
}
func (s *SysResource) GetRoleResourceDetailList(req *dto.GetRoleResourceDetailListReq) (*dto.GetRoleResourceDetailListRes, error) {
	var (
		result        = &dto.GetRoleResourceDetailListRes{}
		db            = s.DB
		log           = s.Log
		resource      = models.SysResourceType{}
		items         = make([]dto.ResourceDetailItem, 0)
		roleResources = make([]*models.SysRoleResource, 0)
	)
	getItemId := func(item dto.ResourceDetailItem) string {
		if item["id"] == nil {
			return ""
		}
		switch item["id"].(type) {
		case int64:
			return strconv.FormatInt(item["id"].(int64), 10)

		case string:
			return item["id"].(string)
		default:
			return ""
		}
	}
	err := db.First(&resource, req.ResourceTypeId).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	getResourceDetailReq := &dto.GetResourceDetailListReq{
		ExtraParams: req.ExtraParams,
		PaginationReq: pkgdto.PaginationReq{
			PageIndex: req.PageIndex,
			PageSize:  req.PageSize,
		},
		Id: req.ResourceTypeId,
	}
	resourceDetail, err := s.GetResourceDetail(getResourceDetailReq)
	if err != nil {
		return nil, err
	}
	// 获取本次查询的所有资源ID
	resourceIds := tools.GetSlice(resourceDetail.Items, func(e dto.ResourceDetailItem) string {
		return getItemId(e)
	})
	err = db.Model(models.SysRoleResource{}).
		Where("resource_type_identify = ? and role_id = ? and resource_id in ?", resource.Identify, req.RoleId, resourceIds).
		Find(&roleResources).Error
	if err != nil {
		err = errors.Wrap(err, "query resource role failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	roleResourceMap := tools.Slice2MapSlice(roleResources, func(e *models.SysRoleResource) string {
		return tools.ToValue(e.ResourceId)
	})
	for _, item := range resourceDetail.Items {
		if roleResource, ok := roleResourceMap[getItemId(item)]; ok {
			if len(roleResource) > 0 {
				permissionTypes := tools.GetSlice(roleResource, func(e *models.SysRoleResource) string {
					return tools.ToValue(e.Action)
				})
				permissionTypes = tools.RemoveDuplication(permissionTypes)
				permissionTypes = tools.SliceFilter(permissionTypes, func(e string) bool {
					return e != ""
				})
				item["permissionTypes"] = permissionTypes
			} else {
				item["permissionTypes"] = []string{}
			}
		}
		items = append(items, item)
	}
	result.Total = resourceDetail.Total
	result.Items = items
	return result, nil
}
func (s *SysResource) UpdateRoleResource(req *dto.UpdateRoleResourceReq) (*dto.UpdateRoleResourceRes, error) {
	var (
		result       = &dto.UpdateRoleResourceRes{}
		db           = s.DB
		log          = s.Log
		resourceType = models.SysResourceType{}
		role         = models.SysRole{}
	)
	// 校验角色是否存在
	err := db.First(&role, req.RoleId).Error
	if err != nil {
		err = errors.Wrap(err, "query role failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	if role.Identify == consts.AdminRoleIdentify {
		return nil, errors.New("系统内置管理员角色不允许修改")
	}
	err = db.First(&resourceType, req.ResourceTypeId).Error
	if err != nil {
		err = errors.Wrap(err, "query resourceType failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	allPermissionCodes := tools.GetSlice(resourceType.PermissionTypes, func(e *models.PermissionType) string {
		return tools.ToValue(e.Code)
	})

	// 校验不存在的权限类型
	for _, authedAllResourcePermissionType := range req.AuthedAllResourcePermissionTypes {
		if !tools.InSlice(authedAllResourcePermissionType, allPermissionCodes) {
			return nil, errors.Errorf("未知的权限类型：%s", authedAllResourcePermissionType)
		}
	}
	for _, changedPermissionType := range req.ChangedPermissionTypes {
		if changedPermissionType == nil {
			continue
		}
		for _, code := range changedPermissionType.PermissionTypes {
			if !tools.InSlice(code, allPermissionCodes) {
				return nil, errors.Errorf("未知的权限类型：%s", code)
			}
			if tools.InSlice(code, req.AuthedAllResourcePermissionTypes) {
				return nil, errors.Errorf("资源（id=%v）已被授权为所有资源权限，无法再设置权限类型：%s",
					tools.ToValue(changedPermissionType.ResourceId), code)
			}
		}
	}
	// 授权
	err = db.Transaction(func(tx *gorm.DB) error {
		// 删除修改过的资源权限信息，重新添加
		// 删除通配符权限，重新添加
		updatedResourceIds := tools.GetSlice(req.ChangedPermissionTypes, func(e *dto.ChangedPermissionType) string {
			return strconv.FormatInt(tools.ToValue(e.ResourceId), 10)
		})
		err = tx.Model(&models.SysRoleResource{}).
			Where("role_id = ? ", req.RoleId).
			Where("resource_type_identify = ? ", resourceType.Identify).
			Scopes(func(tx *gorm.DB) *gorm.DB {
				if len(updatedResourceIds) > 0 {
					tx.Where("resource_id in ?  or resource_id = ? ", updatedResourceIds, auth_engine.Wildcard)
				} else {
					tx.Where("resource_id = ?", auth_engine.Wildcard)
				}
				return tx
			}).Delete(&models.SysRoleResource{}).Error
		if err != nil {
			return errors.Wrap(err, "delete resource role failed")
		}
		var roleResources = make([]*models.SysRoleResource, 0, len(req.AuthedAllResourcePermissionTypes))
		// 添加部分资源授权
		for _, changedPermissionType := range req.ChangedPermissionTypes {
			changedPermissionType.PermissionTypes = tools.RemoveDuplication(changedPermissionType.PermissionTypes)
			for _, permissionType := range changedPermissionType.PermissionTypes {
				roleResources = append(roleResources, &models.SysRoleResource{
					RoleId:               req.RoleId,
					ResourceTypeIdentify: resourceType.Identify,
					ResourceId:           tools.ToPointer(strconv.FormatInt(tools.ToValue(changedPermissionType.ResourceId), 10)),
					Action:               &permissionType,
				})
			}
		}
		// 添加所有资源授权
		for _, permissionType := range req.AuthedAllResourcePermissionTypes {
			roleResources = append(roleResources, &models.SysRoleResource{
				RoleId:               req.RoleId,
				ResourceTypeIdentify: resourceType.Identify,
				ResourceId:           tools.ToPointer(auth_engine.Wildcard),
				Action:               &permissionType,
			})
		}
		err = tx.CreateInBatches(&roleResources, 100).Error
		if err != nil {
			return errors.Wrap(err, "create resource role failed")
		}
		return nil
	})
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}
func (s *SysResource) AuthAllResourceByPermissionCode(db *gorm.DB, roleId int64, resourceIdentify, permissionType string) error {
	var (
		err               error
		existResourceRole = models.SysRoleResource{}
	)

	// 查询是否存在通配符权限，存在则不做操作，不存在则添加
	err = db.Where("resource_type_identify = ? and role_id = ? and resource_id = ? and action=?",
		resourceIdentify, roleId, auth_engine.Wildcard, permissionType).
		Find(&existResourceRole).Error
	if err != nil {
		return errors.Wrap(err, "query resource role failed")
	}
	// 删除之前的非通配符权限，添加新的权限
	err = db.Unscoped().Model(&models.SysRoleResource{}).
		Where("role_id = ? and resource_type_identify = ? and resource_id != ? and action = ?",
			roleId, resourceIdentify, auth_engine.Wildcard, permissionType).
		Delete(&models.SysRoleResource{}).Error
	if err != nil {
		return errors.Wrap(err, "delete resource role failed")
	}
	// 之前存在则不需要操作
	if existResourceRole.Id > 0 {
		return nil
	} else { // 不存在则添加
		err = db.Create(&models.SysRoleResource{
			RoleId:               tools.ToPointer(roleId),
			ResourceTypeIdentify: &resourceIdentify,
			ResourceId:           tools.ToPointer(auth_engine.Wildcard),
			Action:               &permissionType,
		}).Error
		if err != nil {
			return errors.Wrap(err, "add resource role failed")
		}
	}
	return nil
}
func (s *SysResource) AuthPartialResourceByPermissionCode(db *gorm.DB, roleId int64, resourceIdentify string, resourceIds []string, permissionCode string) error {
	roleResources := make([]*models.SysRoleResource, 0)
	// 删除之前的权限，添加新的权限
	err := db.Unscoped().Where("role_id = ? and resource_type_identify = ? and action = ?", roleId, resourceIdentify, permissionCode).
		Delete(&models.SysRoleResource{}).Error
	if err != nil {
		return errors.Wrap(err, "delete resource role failed")
	}
	for _, resourceId := range resourceIds {
		roleResources = append(roleResources, &models.SysRoleResource{
			RoleId:               tools.ToPointer(roleId),
			ResourceTypeIdentify: &resourceIdentify,
			ResourceId:           tools.ToPointer(resourceId),
			Action:               &permissionCode,
		})
	}
	// 添加新权限
	err = db.CreateInBatches(&roleResources, 100).Error
	if err != nil {
		return errors.Wrap(err, "add resource role failed")
	}
	return nil
}
func (s *SysResource) filterPermissionTypes(pts []string, resourceTypes []*models.PermissionType) []string {
	res := make([]string, 0)
	for _, p := range pts {
		for _, pt := range resourceTypes {
			if tools.PtrEqVal(pt.Code, p) {
				res = append(res, p)
			}
		}
	}
	return res
}

func (s *SysResource) GetRoleAuthedResource(req *dto.GetRoleAuthedResourceReq) (*dto.GetRoleAuthedResourceRes, error) {
	var (
		db           = s.DB
		result       = &dto.GetRoleAuthedResourceRes{}
		resourceType = models.SysResourceType{}
		log          = s.Log
		resourceIds  = make([]int64, 0)
	)
	err := db.Scopes(func(tx *gorm.DB) *gorm.DB {
		if tools.ToValue(req.Identify) != "" {
			return tx.Where("identify = ?", req.Identify)
		}
		if tools.ToValue(req.ResourceTypeId) != 0 {
			return tx.Where("id = ?", req.ResourceTypeId)
		}
		return tx
	}).First(&resourceType).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	// 查询是否包含通配符权限
	roleResource := models.SysRoleResource{}
	err = db.Model(&models.SysRoleResource{}).
		Where("role_id = ?", req.RoleId).
		Where("resource_type_identify = ?", req.Identify).
		Where("action = ?", req.PermissionCode).
		Where("resource_id = ?", auth_engine.Wildcard).Find(&roleResource).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	if roleResource.Id > 0 {
		result.AuthAllResource = tools.ToPointer(int64(dto.AuthAllResourceYes))
		return result, nil
	} else {
		result.AuthAllResource = tools.ToPointer(int64(dto.AuthAllResourceNo))
	}

	// 不包含通配符权限
	subQuery := db.Model(models.SysRoleResource{}).Distinct("resource_id").
		Select("resource_id").
		Where("role_id = ?", req.RoleId).
		Where("resource_type_identify = ?", req.Identify).
		Where("action = ?", req.PermissionCode)
	fieldName := "id"
	err = db.Table(tools.ToValue(resourceType.Table)).Select(fieldName).Distinct(fieldName).Scopes(func(tx *gorm.DB) *gorm.DB {
		tx.Where(fieldName + " is not null")
		tx.Where(fieldName + " != '' ")
		// 应用这个字体的筛选条件
		if len(tools.ToValue(resourceType.Filter)) > 0 {
			tx.Where(tools.ToValue(resourceType.Filter))
		}
		tx.Where("id IN (?)", subQuery)
		return tx
	}).Find(&resourceIds).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	result.ResourceIds = resourceIds
	return result, nil
}

func (s *SysResource) GetBusinessResource(req *dto.GetBusinessResourceReq) (*dto.GetBusinessResourceRes, error) {
	var (
		db           = s.DB
		result       = &dto.GetBusinessResourceRes{}
		resourceType = models.SysResourceType{}
		log          = s.Log
		categories   = make([]*models.BusinessCategory, 0)
		data         = make([]*dto.Category, 0)
	)
	err := db.Scopes(func(tx *gorm.DB) *gorm.DB {
		tx.Where("identify = ?", consts.Subcategory)
		return tx
	}).First(&resourceType).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	err = db.Model(models.BusinessCategory{}).
		Preload("Subcategories", func(tx *gorm.DB) *gorm.DB {
			// 值过滤有效的数据
			if tools.ToValue(resourceType.Filter) != "" {
				tx = tx.Where(resourceType.Filter)
			}
			return tx.Order("name asc")
		}).Order("name asc").Find(&categories).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	for _, category := range categories {
		dtoCategory := &dto.Category{
			Id:   tools.ToPointer(category.Id),
			Name: category.Name,
		}
		for _, subCategory := range category.Subcategories {
			dtoSubCategory := &dto.Subcategory{
				Id:   tools.ToPointer(subCategory.Id),
				Name: subCategory.Name,
			}
			dtoCategory.Subcategories = append(dtoCategory.Subcategories, dtoSubCategory)
		}
		data = append(data, dtoCategory)
	}
	result.Data = data
	return result, nil
}
func (s *SysResource) GenerateCategoryIdentify(categoryId int64) string {
	return "category_" + strconv.FormatInt(categoryId, 10)
}
func (s *SysResource) GenerateSubCategoryIdentify(subCategoryId int64) string {
	return "subcategory_" + strconv.FormatInt(subCategoryId, 10)
}

// IsSubCategoryIdentify 判断是否是业务小类ID
func (s *SysResource) IsSubCategoryIdentify(subCategoryIdentify string) bool {
	if !strings.HasPrefix(subCategoryIdentify, "subcategory_") {
		return false
	}
	_, err := strconv.ParseInt(strings.TrimPrefix(subCategoryIdentify, "subcategory_"), 10, 64)
	if err != nil {
		return false
	}
	return true
}

func (s *SysResource) ParseSubCategoryId(subCategoryIdentify string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimPrefix(subCategoryIdentify, "subcategory_"), 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return id, nil
}

func (s *SysResource) RoleResourceAuth(req *dto.RoleResourceAuthReq) (*dto.AuthSubcategoryPermissionRes, error) {
	var (
		result       = &dto.AuthSubcategoryPermissionRes{}
		log          = s.Log
		db           = s.DB
		resourceType = models.SysResourceType{}
	)
	err := db.Scopes(func(tx *gorm.DB) *gorm.DB {
		if tools.ToValue(req.ResourceTypeId) != 0 {
			tx.Where("id = ?", req.ResourceTypeId)
		}
		if tools.ToValue(req.Identify) != "" {
			tx.Where("identify = ?", req.Identify)
		}
		return tx
	}).First(&resourceType).Error
	if err != nil {
		err = errors.Wrap(err, "query resource failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	permissionCodeExist := false
	for _, permissionType := range resourceType.PermissionTypes {
		if tools.PtrEqPtr(req.PermissionCode, permissionType.Code) {
			permissionCodeExist = true
			break
		}
	}
	if !permissionCodeExist {
		err = errors.Errorf("permission code(%v) not exist", tools.ToValue(req.PermissionCode))
		log.Errorf("%v", err)
		return nil, err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if tools.ToValue(req.AuthAllResource) == dto.AuthAllResourceYes {
			err = s.AuthAllResourceByPermissionCode(tx, tools.ToValue(req.RoleId), tools.ToValue(resourceType.Identify), tools.ToValue(req.PermissionCode))
		} else {
			err = s.AuthPartialResourceByPermissionCode(tx, tools.ToValue(req.RoleId), tools.ToValue(resourceType.Identify), req.ResourceIds, tools.ToValue(req.PermissionCode))
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}
