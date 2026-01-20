#!/usr/bin/env node
'use strict';

const fs = require('fs/promises');

/**
 * Convert a document using LibreOffice-as-a-Service
 * @param {Blob|Buffer} fileData - The file data to convert
 * @param {string} filename - Original filename
 * @param {string} format - Target format (e.g., 'pdf', 'txt')
 * @param {object} options - Configuration options
 * @param {string} options.baseUrl - API base URL
 * @param {string} options.apiToken - API authentication token
 * @returns {Promise<Blob|Buffer>} - The converted file
 */
async function convertDocument(fileData, filename, format, options = {}) {
  const baseUrl = options.baseUrl || 'http://127.0.0.1:5227';
  const apiToken = options.apiToken || '';

  const url = new URL(`/api/convert/${format}`, baseUrl);
  url.searchParams.set('filename', filename);

  const headers = {
    'Content-Type': 'application/octet-stream',
  };

  if (apiToken) {
    headers['Authorization'] = `******`;
  }

  const response = await fetch(url.toString(), {
    method: 'POST',
    headers: headers,
    body: fileData,
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Server returned ${response.status}: ${errorText}`);
  }

  // In Node.js, return Buffer; in browser, this would be a Blob
  if (typeof Buffer !== 'undefined') {
    const arrayBuffer = await response.arrayBuffer();
    return Buffer.from(arrayBuffer);
  }
  return await response.blob();
}

// CLI usage (Node.js only)
async function main() {
  const baseUrl = process.env.LAAS_BASE_URL || 'http://127.0.0.1:5227';
  const apiToken = process.env.LAAS_API_TOKEN || '';
  
  const inputFile = 'fixtures/Writing1.docx';
  const outputFile = 'Writing1.pdf';

  try {
    // Read input file
    const fileData = await fs.readFile(inputFile);

    // Convert document
    const converted = await convertDocument(
      fileData,
      'Writing1.docx',
      'pdf',
      { baseUrl, apiToken }
    );

    // Write output file
    await fs.writeFile(outputFile, converted);

    console.log(`Successfully converted ${inputFile} to ${outputFile}`);
  } catch (err) {
    console.error('Error:', err.message);
    process.exit(1);
  }
}

// Run if called directly
if (require.main === module) {
  main();
}

// Export for use as a module (or in browser with bundler)
module.exports = { convertDocument };

/*
 * Browser usage example:
 * 
 * <input type="file" id="fileInput">
 * <button onclick="convert()">Convert to PDF</button>
 * 
 * <script type="module">
 * import { convertDocument } from './example.js';
 * 
 * async function convert() {
 *   const fileInput = document.getElementById('fileInput');
 *   const file = fileInput.files[0];
 *   
 *   const converted = await convertDocument(
 *     file,
 *     file.name,
 *     'pdf',
 *     {
 *       baseUrl: 'http://localhost:5227',
 *       apiToken: 'your-token-here'
 *     }
 *   );
 *   
 *   // Download the converted file
 *   const url = URL.createObjectURL(converted);
 *   const a = document.createElement('a');
 *   a.href = url;
 *   a.download = file.name.replace(/\.[^.]+$/, '.pdf');
 *   a.click();
 *   URL.revokeObjectURL(url);
 * }
 * </script>
 */
