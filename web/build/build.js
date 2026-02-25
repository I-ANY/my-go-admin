import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function getRootPath(...dir) {
  return path.resolve(process.cwd(), ...dir);
}

const runBuild = async () => {
  try {
    const OUTPUT_DIR = 'dist';
    const VERSION = 'version.json';
    const versionJson = {
      version: 'V_' + Math.floor(Math.random() * 10000) + Date.now(),
    };
    console.log(getRootPath(`${OUTPUT_DIR}/${VERSION}`));
    fs.writeFileSync(getRootPath(`${OUTPUT_DIR}/${VERSION}`), JSON.stringify(versionJson));
    console.log(`version file is build successfully!`);
  } catch (error) {
    console.error('version build error:\n' + error);
    process.exit(1);
  }
};

runBuild();
