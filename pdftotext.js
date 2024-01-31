'use strict';

let Os = require('os');
let Path = require('path');
let Fsp = require('fs').promises;
let ChildProcess = require('child_process');

async function exec(cmd, args, opts) {
  return await new Promise(function (resolve, reject) {
    let chunks = [];

    let proc = ChildProcess.spawn(cmd, args, {
      windowsHide: true,
      shell: false,
      timeout: 60 * 1000,
      encoding: 'utf8',
      env: {
        SHELL: process.env.SHELL,
        PATH: process.env.PATH,
        HOME: process.env.HOME,
      },
      cwd: process.cwd(),
    });

    proc.stdout.on('readable', log);
    proc.stderr.on('readable', log);
    proc.on('error', function (err) {
      err.cmd = cmd + ' ' + args.join('');
      err.out = chunks.join('');
      reject(err);
    });
    proc.on('exit', function () {
      resolve(chunks.join());
    });
    if (opts?.stdin) {
      opts.stdin.pipe(proc.stdin);
    }

    function log() {
      let data;
      /*jshint validthis:true*/
      while ((data = this.read())) {
        //console.log(`[spawn: ${cmd}]`, data.toString('utf8'));
        chunks.push(data.toString('utf8'));
      }
    }
  });
}

/**
 * @params String inputPath - full path of input file
 * @params String format - the target format (typically file extension, such as 'pdf')
 * @returns Promise<string> - full path of output file
 */
async function convert(inputPath, options) {
  let cmd = 'pdftotext';
  // ex: /tmp/pdftotext-abc123
  let tmpPrefix = Path.join(Os.tmpdir(), 'pdftotext-');
  let tmpDir = await Fsp.mkdtemp(tmpPrefix);

  let inputExt = Path.extname(inputPath);
  let basename = Path.basename(inputPath, inputExt);
  let outPath = Path.join(tmpDir, `${basename}.txt`);

  let args = ['-eol', 'unix'];
  if (options?.dpi) {
    // TODO
    //
    // -f <int>             : first page to convert
    // -l <int>             : last page to convert
    // -r <fp>              : resolution, in DPI (default is 72)
    // Unit is "dots", so 72 dots is 1" (at 72 DPI)
    // -x <int>             : x-coordinate of the crop area top left corner
    // -y <int>             : y-coordinate of the crop area top left corner
    // -W <int>             : width of crop area in pixels (default is 0)
    // -H <int>             : height of crop area in pixels (default is 0)
  }
  args.push(inputPath);
  args.push(outPath);

  console.info('[pdftotext]', cmd, args.join(' '));
  await exec(cmd, args, { TODO_log: false });
  console.log('[pdftotext] out:', outPath);

  return outPath;
}

module.exports = convert;

if (require.main === module) {
  let input = process.argv[2];
  let outfile = process.argv[3];
  if (!input || !outfile) {
    console.error('Usage: node ./pdftotext.js input.pdf output.txt');
    process.exit(1);
  }

  convert(input)
    .then(async function (out) {
      let Fsp = require('fs').promises;
      await Fsp.rename(out, outfile);
      console.info('Success!', out, outfile);
    })
    .catch(function (err) {
      console.error('Something went wrong:');
      console.error(err);
    });
}
