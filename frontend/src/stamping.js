import { getDocument, PDFWorker } from 'pdfjs-dist';
import { createWorker } from 'tesseract.js';
import { rgb, PDFDocument, StandardFonts } from 'pdf-lib';

export async function stampPDF(data) {
  // Create the service worker for PDF display
  const pdfWorker = new PDFWorker({
    port: new Worker(
      new URL(
        'pdfjs-dist/build/pdf.worker.min.js',
        import.meta.url
      )
    )
  });

  // Load the unstamped PDF document
  const pdfDocument = await getDocument({ data, worker: pdfWorker }).promise;

  // Begin scaling the PDF pages
  const pages = [];
  for (let i = 1; i <= pdfDocument.numPages; i++) {
    const page = await pdfDocument.getPage(i);
    const viewport = page.getViewport({ scale: 4.0 });

    // The `canvas` element is used to render the PDF page
    const canvas = document.createElement('canvas');
    canvas.height = viewport.height;
    canvas.width = viewport.width;

    // Perform the scaled rendering and add the Base64 encoded page to the array
    await page.render({ canvasContext: canvas.getContext('2d'), viewport }).promise;
    pages.push(canvas.toDataURL('image/jpeg'));
  }

  // Create the service worker for OCR
  const tesseractWorker = await createWorker('eng', 1, {
    cachePath: './assets/resources'
  });

  // Perform OCR on the scaled PDF pages
  const response = [];
  for (const page of pages) {
    const {
      data: { pdf }
    } = await tesseractWorker.recognize(page, { pdfTitle: 'Stamped PDF' }, { pdf: true });

    response.push(Uint8Array.from(pdf));
  }

  // Load the unstamped, scaled, OCR'd PDF document for stamping
  const pdfDoc = await PDFDocument.load(await new Blob(response).arrayBuffer());
  const helveticaFont = await pdfDoc.embedFont(StandardFonts.Helvetica);

  // Stamp the first page of the PDF document and scale the pages back down
  let firstPage = true;
  for (const page of pdfDoc.getPages()) {
    // Scale the page back down to its original size
    page.scale(0.125, 0.125);

      // Stamp the first page only
      if (firstPage) {
        firstPage = false;
        const size = page.getSize();

        // Perform the stamping
        page.drawText(`FILED ${new Date().toLocaleDateString()} NYS Department of State`, {
          x: 200,
          y: size.height * 7.7,
          size: 50,
          font: helveticaFont,
          color: rgb(0.95, 0.1, 0.1)
        });
      }
    }

    // Save the stamped PDF document
    const pdfBytes = await pdfDoc.save();

    // Terminate the service workers
    await tesseractWorker.terminate();
    pdfWorker.destroy();

    // Return the stamped PDF document as an ArrayBuffer
    return pdfBytes;
}
