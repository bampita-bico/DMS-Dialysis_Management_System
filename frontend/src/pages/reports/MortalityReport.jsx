import { useMemo, useState } from 'react';
import useOfflineData from '../../hooks/useOfflineData';
import { authService } from '../../services/auth';
import { escapeHtml, htmlTable, printHtml } from '../../utils/print';

const deathSettingLabels = {
  during_dialysis: 'During dialysis',
  hospital: 'Hospital',
  home: 'Home',
  transit: 'Transit',
  other: 'Other',
};

export default function MortalityReport() {
  const [reportMonth, setReportMonth] = useState(() => new Date().toISOString().slice(0, 7));
  const { data: mortalityRecords, loading, error, refresh } = useOfflineData('mortality_records');
  const { data: patients } = useOfflineData('patients');
  const { data: sessions } = useOfflineData('dialysis_sessions');
  const user = authService.getCurrentUser();

  const report = useMemo(() => buildMonthlyReport({
    reportMonth,
    mortalityRecords: mortalityRecords || [],
    patients: patients || [],
    sessions: sessions || [],
    hospitalName: user?.hospital_name || 'Dialysis Unit',
  }), [reportMonth, mortalityRecords, patients, sessions, user?.hospital_name]);

  const handlePrint = () => {
    printHtml({
      title: 'Monthly Mortality Report',
      subtitle: `${report.hospitalName} | ${report.monthLabel}`,
      body: buildPrintBody(report),
      footer: 'Prepared from DMS dialysis-unit records. Reconcile against the official MoH HMIS/DHIS2 submission workflow before final submission.',
    });
  };

  const handleExportCsv = () => {
    const csv = buildCsv(report.rows);
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `monthly-mortality-report-${reportMonth}.csv`;
    link.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div>
              <h1 className="text-3xl font-serif font-bold text-gray-900 tracking-tight">
                Monthly Mortality Report
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                MMR for dialysis-unit death review, ministry reconciliation, and DHIS2/HMIS preparation.
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <label className="flex items-center gap-2 text-sm font-medium text-gray-700">
                Month
                <input
                  type="month"
                  value={reportMonth}
                  onChange={(event) => setReportMonth(event.target.value)}
                  className="rounded-lg border border-gray-300 px-3 py-2"
                />
              </label>
              <button
                onClick={() => refresh()}
                className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Refresh
              </button>
              <button
                onClick={handleExportCsv}
                disabled={report.rows.length === 0}
                className="rounded-lg border border-sky-200 bg-sky-50 px-4 py-2 text-sm font-medium text-sky-700 hover:bg-sky-100 disabled:opacity-50"
              >
                Export CSV
              </button>
              <button
                onClick={handlePrint}
                className="rounded-lg bg-sky-600 px-4 py-2 text-sm font-medium text-white hover:bg-sky-700"
              >
                Print MMR
              </button>
            </div>
          </div>
        </div>
      </div>

      <main className="max-w-7xl mx-auto px-6 py-8 space-y-6">
        {error && (
          <div className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        <section className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-5">
          <MetricCard label="Deaths" value={report.totalDeaths} detail={report.monthLabel} tone={report.totalDeaths ? 'red' : 'green'} />
          <MetricCard label="Session-related" value={report.sessionRelatedDeaths} detail="During or linked to dialysis" tone={report.sessionRelatedDeaths ? 'red' : 'green'} />
          <MetricCard label="Certified" value={report.certifiedDeaths} detail={`${report.pendingCertification} pending`} tone={report.pendingCertification ? 'yellow' : 'green'} />
          <MetricCard label="Active patients" value={report.activePatients} detail="Current local census" tone="blue" />
          <MetricCard label="Death rate" value={`${report.deathRatePer100} / 100`} detail="Monthly active-patient rate" tone={report.totalDeaths ? 'yellow' : 'green'} />
        </section>

        <div className="grid grid-cols-1 gap-6 xl:grid-cols-3">
          <Panel title="Deaths By Setting">
            <SummaryList rows={report.settingSummary} />
          </Panel>
          <Panel title="Common Causes">
            <SummaryList rows={report.causeSummary} emptyText="No deaths in selected month." />
          </Panel>
          <Panel title="Submission Checklist">
            <ul className="space-y-2 text-sm text-gray-700">
              <li>Death record present for each deceased dialysis patient.</li>
              <li>Primary cause and ICD-10 code reviewed by clinician.</li>
              <li>Death certificate number entered where available.</li>
              <li>Session-related deaths flagged for unit quality review.</li>
              <li>MMR reconciled with official MoH HMIS/DHIS2 submission.</li>
            </ul>
          </Panel>
        </div>

        <Panel title={`Mortality Register - ${report.monthLabel}`}>
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="h-12 w-12 animate-spin rounded-full border-b-2 border-sky-600"></div>
            </div>
          ) : report.rows.length === 0 ? (
            <p className="text-sm text-gray-600">No mortality records found for this month.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200 text-sm">
                <thead>
                  <tr className="text-left text-xs font-semibold uppercase tracking-wide text-gray-500">
                    <th className="py-3 pr-4">Date</th>
                    <th className="py-3 pr-4">Patient</th>
                    <th className="py-3 pr-4">MRN</th>
                    <th className="py-3 pr-4">Age/Sex</th>
                    <th className="py-3 pr-4">Setting</th>
                    <th className="py-3 pr-4">Session-related</th>
                    <th className="py-3 pr-4">Cause</th>
                    <th className="py-3 pr-4">ICD-10</th>
                    <th className="py-3 pr-4">Certificate</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {report.rows.map(row => (
                    <tr key={row.id}>
                      <td className="py-3 pr-4 text-gray-700">{row.dateOfDeath}</td>
                      <td className="py-3 pr-4 font-medium text-gray-900">{row.patientName}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.mrn}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.ageSex}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.deathSetting}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.sessionRelated}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.primaryCause}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.icd10Code}</td>
                      <td className="py-3 pr-4 text-gray-700">{row.certificate}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Panel>
      </main>
    </div>
  );
}

function buildMonthlyReport({ reportMonth, mortalityRecords, patients, sessions, hospitalName }) {
  const [year, month] = reportMonth.split('-').map(Number);
  const start = new Date(year, month - 1, 1);
  const end = new Date(year, month, 0, 23, 59, 59, 999);
  const patientById = new Map(patients.map(patient => [String(patient.id), patient]));
  const sessionById = new Map(sessions.map(session => [String(session.id), session]));
  const activePatients = patients.filter(patient => patient.is_active !== false && !patient.deleted_at).length;

  const rows = mortalityRecords
    .filter(record => {
      const time = toTime(record.date_of_death);
      return time >= start.getTime() && time <= end.getTime() && !record.deleted_at;
    })
    .sort((a, b) => toTime(a.date_of_death) - toTime(b.date_of_death))
    .map(record => {
      const patient = patientById.get(String(record.patient_id));
      const session = record.session_id ? sessionById.get(String(record.session_id)) : null;
      return {
        id: record.id,
        dateOfDeath: formatDate(record.date_of_death),
        patientName: patient?.full_name || record.patient_name || 'Unknown patient',
        mrn: patient?.mrn || 'N/A',
        ageSex: `${calculateAge(patient?.date_of_birth)} / ${humanize(patient?.sex || 'N/A')}`,
        deathSetting: deathSettingLabels[record.death_setting] || humanize(record.death_setting || 'other'),
        sessionRelated: record.session_related ? 'Yes' : 'No',
        primaryCause: record.primary_cause_of_death || 'Not recorded',
        contributingFactors: record.contributing_factors || 'None recorded',
        icd10Code: record.icd10_code || 'N/A',
        certificate: record.death_certificate_number || (record.certified_at ? 'Certified' : 'Pending'),
        certified: Boolean(record.death_certificate_number || record.certified_at),
        sessionDate: formatDate(session?.scheduled_date),
        notes: record.notes || '',
      };
    });

  const totalDeaths = rows.length;
  const sessionRelatedDeaths = rows.filter(row => row.sessionRelated === 'Yes').length;
  const certifiedDeaths = rows.filter(row => row.certified).length;
  const pendingCertification = totalDeaths - certifiedDeaths;

  return {
    hospitalName,
    reportMonth,
    monthLabel: start.toLocaleDateString(undefined, { month: 'long', year: 'numeric' }),
    periodLabel: `${formatDate(start)} to ${formatDate(end)}`,
    rows,
    totalDeaths,
    sessionRelatedDeaths,
    certifiedDeaths,
    pendingCertification,
    activePatients,
    deathRatePer100: activePatients ? ((totalDeaths / activePatients) * 100).toFixed(1) : '0.0',
    settingSummary: summarize(rows, 'deathSetting'),
    causeSummary: summarize(rows, 'primaryCause'),
  };
}

function buildPrintBody(report) {
  const metricBoxes = [
    ['Reporting facility', report.hospitalName],
    ['Reporting month', report.monthLabel],
    ['Period', report.periodLabel],
    ['Total deaths', report.totalDeaths],
    ['Session-related deaths', report.sessionRelatedDeaths],
    ['Certified deaths', report.certifiedDeaths],
    ['Pending certification', report.pendingCertification],
    ['Monthly death rate', `${report.deathRatePer100} per 100 active patients`],
  ].map(([label, value]) => (
    `<div class="box"><span class="label">${escapeHtml(label)}</span><span class="value">${escapeHtml(value)}</span></div>`
  )).join('');

  return `
    <div class="meta">${metricBoxes}</div>
    <h2>Deaths By Setting</h2>
    ${htmlTable(['Setting', 'Deaths'], report.settingSummary.map(item => [item.label, item.count]))}
    <h2>Common Causes</h2>
    ${htmlTable(['Cause', 'Deaths'], report.causeSummary.map(item => [item.label, item.count]))}
    <h2>Mortality Register</h2>
    ${htmlTable(
      ['Date', 'Patient', 'MRN', 'Age/Sex', 'Setting', 'Session-related', 'Cause', 'ICD-10', 'Certificate'],
      report.rows.map(row => [row.dateOfDeath, row.patientName, row.mrn, row.ageSex, row.deathSetting, row.sessionRelated, row.primaryCause, row.icd10Code, row.certificate])
    )}
    <div class="signatures">
      <div class="line">Prepared by / Records officer</div>
      <div class="line">Reviewed by / Clinical lead</div>
    </div>
  `;
}

function buildCsv(rows) {
  const headers = ['Date', 'Patient', 'MRN', 'Age/Sex', 'Setting', 'Session Related', 'Primary Cause', 'Contributing Factors', 'ICD-10', 'Certificate', 'Session Date', 'Notes'];
  const dataRows = rows.map(row => [
    row.dateOfDeath,
    row.patientName,
    row.mrn,
    row.ageSex,
    row.deathSetting,
    row.sessionRelated,
    row.primaryCause,
    row.contributingFactors,
    row.icd10Code,
    row.certificate,
    row.sessionDate,
    row.notes,
  ]);
  return [headers, ...dataRows].map(row => row.map(csvCell).join(',')).join('\n');
}

function csvCell(value) {
  return `"${String(value ?? '').replaceAll('"', '""')}"`;
}

function summarize(rows, key) {
  const counts = rows.reduce((acc, row) => {
    const label = row[key] || 'Not recorded';
    acc.set(label, (acc.get(label) || 0) + 1);
    return acc;
  }, new Map());

  return Array.from(counts.entries())
    .map(([label, count]) => ({ label, count }))
    .sort((a, b) => b.count - a.count || a.label.localeCompare(b.label));
}

function Panel({ title, children }) {
  return (
    <section className="rounded-lg border border-gray-200 bg-white p-5">
      <h2 className="mb-4 text-lg font-semibold text-gray-900">{title}</h2>
      {children}
    </section>
  );
}

function MetricCard({ label, value, detail, tone = 'gray' }) {
  return (
    <div className={`rounded-lg border p-4 ${toneClasses(tone).soft}`}>
      <p className="text-xs font-semibold uppercase tracking-wide text-gray-500">{label}</p>
      <p className={`mt-2 text-2xl font-bold ${toneClasses(tone).text}`}>{value}</p>
      <p className="mt-1 text-sm text-gray-700">{detail}</p>
    </div>
  );
}

function SummaryList({ rows, emptyText = 'No records found.' }) {
  if (rows.length === 0) return <p className="text-sm text-gray-600">{emptyText}</p>;
  return (
    <div className="space-y-2">
      {rows.map(row => (
        <div key={row.label} className="flex items-center justify-between rounded-lg border border-gray-200 bg-gray-50 px-3 py-2">
          <span className="text-sm font-medium text-gray-800">{row.label}</span>
          <span className="text-sm font-bold text-gray-900">{row.count}</span>
        </div>
      ))}
    </div>
  );
}

function toTime(value) {
  if (!value) return 0;
  const time = new Date(value).getTime();
  return Number.isNaN(time) ? 0 : time;
}

function formatDate(value) {
  if (!value) return 'Not recorded';
  const date = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleDateString();
}

function calculateAge(dateOfBirth) {
  if (!dateOfBirth) return 'N/A';
  const birth = new Date(dateOfBirth);
  if (Number.isNaN(birth.getTime())) return 'N/A';
  return Math.floor((new Date() - birth) / 31557600000);
}

function humanize(value) {
  if (!value) return 'Not recorded';
  return String(value)
    .replaceAll('_', ' ')
    .replaceAll('-', ' ')
    .replace(/\s+/g, ' ')
    .trim()
    .replace(/\b\w/g, char => char.toUpperCase());
}

function toneClasses(tone) {
  const classes = {
    green: { soft: 'bg-emerald-50 border-emerald-200', text: 'text-emerald-800' },
    red: { soft: 'bg-red-50 border-red-200', text: 'text-red-800' },
    yellow: { soft: 'bg-amber-50 border-amber-200', text: 'text-amber-800' },
    blue: { soft: 'bg-sky-50 border-sky-200', text: 'text-sky-800' },
    gray: { soft: 'bg-gray-50 border-gray-200', text: 'text-gray-800' },
  };
  return classes[tone] || classes.gray;
}
