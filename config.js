'use strict';

// Note: You MUST load any ENVs BEFORE this file is required. Example:
// require('dotenv').config({ path: ".env" })

let config = module.exports;

config.NODE_ENV = process.env.NODE_ENV || 'development';
config.PORT = process.env.PORT || '5227';

config.API_TOKEN = process.env.LAAS_API_TOKEN || process.env.API_TOKEN || '';
