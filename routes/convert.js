'use strict';

let tmpdir = require('os').tmpdir();
let Path = require('path');
let Fs = require('fs');
let Fsp = require('fs').promises;

let Config = require('../config.js');
let Convert = require('../');

module.exports = function (fastify, opts, done) {
  fastify.addContentTypeParser('*', function (request, payload, done2) {
    // skip parsing so that the handler has access to the stream
    //request.raw.pause();
    console.log('[DEBUG] rando type?');
    console.log(request.raw.rawHeaders);
    done2(null, undefined);
  });

  async function receive(req, name, format = 'pdf') {
    let tmpPrefix = Path.join(tmpdir, 'laas-source-');
    let dirname = await Fsp.mkdtemp(tmpPrefix);
    console.info('[receive] incoming tmpdir:', dirname);
    let dst = Path.join(dirname, name);
    let stream = Fs.createWriteStream(dst);

    let originalPath = await new Promise(function (resolve, reject) {
      req.pipe(stream);
      req.on('readable', function () {
        let chunk;
        while ((chunk = req.read())) {
          console.log('[DEBUG] chunk.length', chunk.length);
        }
      });
      req.on('error', reject);
      req.on('end', function () {
        // remember: close() for files, end() for network streams
        stream.close();
        stream.on('close', function () {
          resolve(dst);
        });
      });
    });

    return originalPath;
  }

  // POST /api/convert/:name (ex: report.docx)
  fastify.post('/:format', async function (request, reply) {
    {
      let token = (request.raw.headers.authorization || '').split(' ')[1];
      let tokenMatches = secureCompare(token, Config.API_TOKEN);
      if (!tokenMatches) {
        let skipToken = '*' === Config.API_TOKEN;
        if (!skipToken) {
          reply.code(401);
          return { success: false, error: 'UNAUTHORIZED' };
        }
      }
    }

    // target format
    //@ts-ignore - not sure why it's not picking up the type
    let format = request.params.format;

    //@ts-ignore - not sure why it's not picking up the type
    let filename = request.query.filename;

    // content-disposition filename hint
    filename = Path.basename(filename);
    if (!filename) {
      throw new Error("BAD_REQUEST: 'filename' should be the name of the source file");
    }

    let originalPath = await receive(request.raw, filename, format);
    let convertedPath = await Convert.convert(originalPath, format);
    let stream = Fs.createReadStream(convertedPath);
    stream.on('end', async function () {
      await Fsp.unlink(convertedPath).catch(function (err) {
        console.error(`Error: failed to remove ${convertedPath}`);
      });
    });

    let suggestedName = Path.basename(convertedPath);
    suggestedName = suggestedName.replace(/"/g, '\\"');
    reply.header('Content-Disposition', `attachment; filename="${suggestedName}"`);
    reply.send(stream);
  });

  function secureCompare(a, b) {
    if (!a && !b) {
      throw new Error('[secure compare] reference string should not be empty');
    }

    if (a.length !== b.length) {
      return false;
    }

    let Crypto = require('crypto');
    return Crypto.timingSafeEqual(Buffer.from(a), Buffer.from(b));
  }

  done();
};
