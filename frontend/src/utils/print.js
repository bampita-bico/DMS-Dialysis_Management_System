export function escapeHtml(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}

export function printHtml({ title, subtitle = '', body, footer = '' }) {
  const printWindow = window.open('', '_blank');
  if (!printWindow) {
    window.print();
    return;
  }

  printWindow.document.write(`<!doctype html>
<html>
  <head>
    <title>${escapeHtml(title)}</title>
    <style>
      * { box-sizing: border-box; }
      body { color: #111827; font-family: Arial, sans-serif; margin: 32px; }
      h1 { font-size: 24px; margin: 0 0 4px; }
      h2 { border-bottom: 1px solid #d1d5db; font-size: 16px; margin: 24px 0 10px; padding-bottom: 6px; }
      h3 { font-size: 14px; margin: 18px 0 8px; }
      p { margin: 4px 0; }
      table { border-collapse: collapse; margin-top: 8px; width: 100%; }
      th, td { border: 1px solid #d1d5db; font-size: 11px; padding: 6px; text-align: left; vertical-align: top; }
      th { background: #f3f4f6; font-weight: 700; }
      .subtitle { color: #4b5563; font-size: 12px; margin-bottom: 16px; }
      .meta { display: grid; gap: 8px; grid-template-columns: repeat(4, 1fr); margin: 16px 0; }
      .box { border: 1px solid #d1d5db; padding: 10px; }
      .label { color: #6b7280; display: block; font-size: 10px; font-weight: 700; text-transform: uppercase; }
      .value { display: block; font-size: 16px; font-weight: 700; margin-top: 4px; }
      .muted { color: #6b7280; font-size: 11px; }
      .signatures { display: grid; gap: 24px; grid-template-columns: 1fr 1fr; margin-top: 32px; }
      .line { border-top: 1px solid #111827; padding-top: 6px; }
      @page { margin: 18mm; }
      @media print {
        body { margin: 0; }
        button { display: none; }
      }
    </style>
  </head>
  <body>
    <h1>${escapeHtml(title)}</h1>
    ${subtitle ? `<p class="subtitle">${escapeHtml(subtitle)}</p>` : ''}
    ${body}
    ${footer ? `<p class="muted">${escapeHtml(footer)}</p>` : ''}
  </body>
</html>`);

  printWindow.document.close();
  printWindow.focus();
  window.setTimeout(() => printWindow.print(), 150);
}

export function keyValueTable(rows) {
  return `<table><tbody>${rows
    .filter(([, value]) => value !== undefined && value !== null && value !== '')
    .map(([label, value]) => `<tr><th>${escapeHtml(label)}</th><td>${escapeHtml(value)}</td></tr>`)
    .join('')}</tbody></table>`;
}

export function htmlTable(headers, rows) {
  return `<table>
    <thead><tr>${headers.map(header => `<th>${escapeHtml(header)}</th>`).join('')}</tr></thead>
    <tbody>${rows.map(row => `<tr>${row.map(value => `<td>${escapeHtml(value)}</td>`).join('')}</tr>`).join('')}</tbody>
  </table>`;
}
