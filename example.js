#!/usr/bin/env node
'use strict';

const fs = require('fs');
const http = require('http');

const BASE_URL = process.env.LAAS_BASE_URL || 'http://127.0.0.1:5227';
const API_TOKEN = process.env.LAAS_API_TOKEN || '';

// Read the input file
const inputFile = 'fixtures/Writing1.docx';
const outputFile = 'Writing1.pdf';

fs.readFile(inputFile, (err, data) => {
  if (err) {
    console.error('Error reading file:', err);
    process.exit(1);
  }

  const url = new URL('/api/convert/pdf', BASE_URL);
  url.searchParams.set('filename', 'Writing1.docx');

  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/octet-stream',
      'Content-Length': data.length,
    },
  };

  if (API_TOKEN) {
    options.headers['Authorization'] = `******`;
  }

  const req = http.request(url, options, (res) => {
    if (res.statusCode !== 200) {
      console.error(`Error: Server returned ${res.statusCode}`);
      res.pipe(process.stderr);
      process.exit(1);
    }

    const writeStream = fs.createWriteStream(outputFile);
    res.pipe(writeStream);

    writeStream.on('finish', () => {
      console.log(`Successfully converted ${inputFile} to ${outputFile}`);
    });

    writeStream.on('error', (err) => {
      console.error('Error writing file:', err);
      process.exit(1);
    });
  });

  req.on('error', (err) => {
    console.error('Request error:', err);
    process.exit(1);
  });

  req.write(data);
  req.end();
});
