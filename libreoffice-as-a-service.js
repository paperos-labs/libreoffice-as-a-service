'use strict';

let Convert = module.exports;

let Path = require('path');

let sofficeConvert = require('./soffice.js');
let pdftotextConvert = require('./pdftotext.js');

/**
 * @param {String} originalPath
 * @param {String} format
 */
Convert.convert = async function (originalPath, format) {
  let fromExt = Path.extname(originalPath);
  fromExt = fromExt.toLowerCase();

  let convertedPath;
  let isPdfToTxt = fromExt === '.pdf' && format === 'txt';
  if (isPdfToTxt) {
    convertedPath = await pdftotextConvert(originalPath);
    return convertedPath;
  }

  convertedPath = await sofficeConvert(originalPath, format);
  return convertedPath;
};
