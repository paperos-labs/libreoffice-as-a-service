'use strict';

let pkg = require('../package.json');

module.exports = function (fastify, opts, done) {
  // GET /api/versions
  fastify.get('/', async function (request, reply) {
    let versions = {
      api: pkg.version,
      node: process.versions.node,
      icu: process.versions.icu,
      tz: process.versions.tz,
    };

    reply.send(versions);
  });

  done();
};
