'use client';

import { useEffect, useState } from 'react';

const API = 'http://localhost:8080';

interface ConversationMetrics {
  success_rate: number;
  intent_accuracy: number;
  avg_steps_completed: number;
  multi_turn_count: number;
  context_retention: number;
  coherence_score: number;
  completeness_score: number;
  sentiment_score: number;
  confidence_score: number;
  avg_response_length: number;
  total_tests: number;
  passed_tests: number;
  avg_latency_ms: number;
}

interface TestResult {
  id: number;
  scenario: string;
  passed: boolean;
  latency_ms: number;
  error: string;
  variant: string;
  created_at: string;
}

interface ABTest {
  id: number;
  test_name: string;
  variant_a: string;
  variant_b: string;
  winner: string;
  summary: string;
  created_at: string;
}

interface Transcription {
  id: number;
  call_sid: string;
  audio_path: string;
  text: string;
  confidence: number;
  duration_seconds: number;
  created_at: string;
}

type Tab = 'metrics' | 'results' | 'ab-tests' | 'transcriptions' | 'health';

export default function Home() {
  const [tab, setTab] = useState<Tab>('metrics');
  const [metrics, setMetrics] = useState<ConversationMetrics | null>(null);
  const [results, setResults] = useState<TestResult[]>([]);
  const [abTests, setABTests] = useState<ABTest[]>([]);
  const [transcriptions, setTranscriptions] = useState<Transcription[]>([]);
  const [health, setHealth] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const safeFetch = (url: string) =>
      fetch(url).then(r => r.ok ? r.json() : null).catch(() => null);

    Promise.all([
      safeFetch(`${API}/api/metrics`),
      safeFetch(`${API}/api/results?limit=20`),
      safeFetch(`${API}/api/ab-tests`),
      safeFetch(`${API}/api/transcriptions`),
      safeFetch(`${API}/api/health/status`),
    ]).then(([m, r, ab, t, h]) => {
        setMetrics(m);
        setResults(r || []);
        setABTests(ab || []);
        setTranscriptions(t || []);
        setHealth(h);
        setLoading(false);
      });
  }, []);

  if (loading) return <Center>Loading...</Center>;

  const tabs: { id: Tab; label: string }[] = [
    { id: 'metrics', label: '📊 Metrics' },
    { id: 'results', label: '🧪 Test Results' },
    { id: 'ab-tests', label: '🔬 A/B Tests' },
    { id: 'transcriptions', label: '🎙️ Transcriptions' },
    { id: 'health', label: '❤️ Health' },
  ];

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-1">🐦 MockingJay</h1>
        <p className="text-gray-500 mb-6">Voice AI Testing Platform</p>

        {/* Tabs */}
        <div className="flex gap-2 mb-6 border-b border-gray-200">
          {tabs.map(t => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`px-4 py-2 text-sm font-medium rounded-t-lg transition-colors ${
                tab === t.id
                  ? 'bg-white border border-b-white border-gray-200 text-blue-600 -mb-px'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>

        {tab === 'metrics' && metrics && <MetricsTab metrics={metrics} />}
        {tab === 'results' && <ResultsTab results={results} />}
        {tab === 'ab-tests' && <ABTestsTab tests={abTests} />}
        {tab === 'transcriptions' && <TranscriptionsTab transcriptions={transcriptions} />}
        {tab === 'health' && <HealthTab health={health} />}
      </div>
    </div>
  );
}

function MetricsTab({ metrics }: { metrics: ConversationMetrics }) {
  return (
    <div className="space-y-6">
      {/* Summary */}
      <div className="grid grid-cols-3 gap-4">
        <MetricCard label="Total Tests" value={String(metrics.total_tests)} color="blue" />
        <MetricCard label="Passed" value={String(metrics.passed_tests)} color="green" />
        <MetricCard label="Avg Latency" value={`${metrics.avg_latency_ms?.toFixed(0) ?? 0}ms`} color="purple" />
      </div>

      <Section title="💬 Conversation Intelligence">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <MetricCard label="Success Rate" value={`${metrics.success_rate.toFixed(1)}%`} color="green" />
          <MetricCard label="Intent Accuracy" value={`${metrics.intent_accuracy.toFixed(1)}%`} color="blue" />
          <MetricCard label="Avg Steps" value={metrics.avg_steps_completed.toFixed(1)} color="purple" />
          <MetricCard label="Multi-turn" value={String(metrics.multi_turn_count)} color="indigo" />
        </div>
      </Section>

      <Section title="🔄 Multi-turn Dialogue">
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          <MetricCard label="Context Retention" value={`${metrics.context_retention.toFixed(1)}%`} color="green" />
          <MetricCard label="Coherence Score" value={`${metrics.coherence_score.toFixed(1)}%`} color="blue" />
          <MetricCard label="Multi-turn Count" value={String(metrics.multi_turn_count)} color="purple" />
        </div>
      </Section>

      <Section title="✨ Response Quality">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <MetricCard label="Completeness" value={`${metrics.completeness_score.toFixed(1)}%`} color="green" />
          <MetricCard label="Sentiment" value={`${metrics.sentiment_score.toFixed(1)}%`} color="yellow" />
          <MetricCard label="Confidence" value={`${metrics.confidence_score.toFixed(1)}%`} color="blue" />
          <MetricCard label="Avg Length" value={`${metrics.avg_response_length.toFixed(0)} chars`} color="purple" />
        </div>
      </Section>
    </div>
  );
}

function ResultsTab({ results }: { results: TestResult[] }) {
  if (results.length === 0) {
    return <Empty>No test results yet. Run <code>mockingjay run --api-url http://localhost:8080</code></Empty>;
  }
  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <table className="w-full text-sm">
        <thead className="bg-gray-50 text-gray-600 uppercase text-xs">
          <tr>
            <th className="px-4 py-3 text-left">Scenario</th>
            <th className="px-4 py-3 text-left">Status</th>
            <th className="px-4 py-3 text-left">Latency</th>
            <th className="px-4 py-3 text-left">Variant</th>
            <th className="px-4 py-3 text-left">Time</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {results.map(r => (
            <tr key={r.id} className="hover:bg-gray-50">
              <td className="px-4 py-3 font-medium">{r.scenario}</td>
              <td className="px-4 py-3">
                <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${r.passed ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                  {r.passed ? '✓ PASS' : '✗ FAIL'}
                </span>
              </td>
              <td className="px-4 py-3 text-gray-600">{r.latency_ms}ms</td>
              <td className="px-4 py-3 text-gray-500">{r.variant || '—'}</td>
              <td className="px-4 py-3 text-gray-400">{new Date(r.created_at).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function ABTestsTab({ tests }: { tests: ABTest[] }) {
  if (tests.length === 0) {
    return <Empty>No A/B tests yet. Run <code>mockingjay ab</code></Empty>;
  }
  return (
    <div className="space-y-4">
      {tests.map(t => (
        <div key={t.id} className="bg-white rounded-lg shadow p-5">
          <div className="flex items-center justify-between mb-3">
            <h3 className="font-semibold text-gray-900">{t.test_name}</h3>
            <span className="text-xs text-gray-400">{new Date(t.created_at).toLocaleString()}</span>
          </div>
          <div className="flex gap-4 text-sm mb-2">
            <span className="text-gray-600">A: <strong>{t.variant_a}</strong></span>
            <span className="text-gray-400">vs</span>
            <span className="text-gray-600">B: <strong>{t.variant_b}</strong></span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-sm text-gray-500">Winner:</span>
            <span className="px-2 py-0.5 bg-yellow-100 text-yellow-800 rounded-full text-xs font-medium">
              🏆 {t.winner}
            </span>
          </div>
          {t.summary && <p className="mt-2 text-sm text-gray-500">{t.summary}</p>}
        </div>
      ))}
    </div>
  );
}

function TranscriptionsTab({ transcriptions }: { transcriptions: Transcription[] }) {
  if (transcriptions.length === 0) {
    return <Empty>No transcriptions yet. Run <code>mockingjay transcribe --file recording.wav</code></Empty>;
  }
  return (
    <div className="space-y-4">
      {transcriptions.map(t => (
        <div key={t.id} className="bg-white rounded-lg shadow p-5">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs text-gray-400 font-mono">{t.audio_path}</span>
            <span className="text-xs text-gray-400">{new Date(t.created_at).toLocaleString()}</span>
          </div>
          <p className="text-gray-800 mb-3">{t.text}</p>
          <div className="flex gap-4 text-xs text-gray-500">
            <span>Confidence: {(t.confidence * 100).toFixed(1)}%</span>
            <span>Duration: {t.duration_seconds.toFixed(1)}s</span>
            {t.call_sid && <span>Call: {t.call_sid}</span>}
          </div>
        </div>
      ))}
    </div>
  );
}

function HealthTab({ health }: { health: any }) {
  if (!health) return <Empty>No health data yet. Run <code>mockingjay run --api-url http://localhost:8080</code></Empty>;

  const statusColors: Record<string, string> = {
    healthy: 'bg-green-100 text-green-800 border-green-300',
    degraded: 'bg-yellow-100 text-yellow-800 border-yellow-300',
    unhealthy: 'bg-red-100 text-red-800 border-red-300',
  };
  const statusColor = statusColors[health.status] ?? statusColors.healthy;

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-lg shadow p-6 flex items-center gap-6">
        <div className={`border-2 rounded-lg px-6 py-4 text-center ${statusColor}`}>
          <div className="text-xs font-medium opacity-75 mb-1">Status</div>
          <div className="text-2xl font-bold capitalize">{health.status}</div>
        </div>
        <div>
          <div className="text-3xl font-bold text-gray-900">{health.pass_rate_24h?.toFixed(1)}%</div>
          <div className="text-sm text-gray-500">Pass rate (last 24h)</div>
          <div className="text-xs text-gray-400 mt-1">{health.passed_24h} passed / {health.total_24h} total</div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="px-5 py-3 border-b border-gray-100 font-medium text-gray-700">Recent Runs</div>
        {health.recent_runs?.length === 0
          ? <div className="p-5 text-gray-400 text-sm">No recent runs</div>
          : <table className="w-full text-sm">
              <thead className="bg-gray-50 text-gray-500 uppercase text-xs">
                <tr>
                  <th className="px-4 py-2 text-left">Scenario</th>
                  <th className="px-4 py-2 text-left">Status</th>
                  <th className="px-4 py-2 text-left">Latency</th>
                  <th className="px-4 py-2 text-left">Time</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {health.recent_runs?.map((r: any, i: number) => (
                  <tr key={i} className="hover:bg-gray-50">
                    <td className="px-4 py-2 font-medium">{r.scenario}</td>
                    <td className="px-4 py-2">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${r.passed ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                        {r.passed ? '✓ PASS' : '✗ FAIL'}
                      </span>
                    </td>
                    <td className="px-4 py-2 text-gray-500">{r.latency_ms}ms</td>
                    <td className="px-4 py-2 text-gray-400">{r.created_at}</td>
                  </tr>
                ))}
              </tbody>
            </table>
        }
      </div>
    </div>
  );
}

// Shared components
function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-semibold mb-4">{title}</h2>
      {children}
    </div>
  );
}

function MetricCard({ label, value, color }: { label: string; value: string; color: string }) {
  const colors: Record<string, string> = {
    green: 'bg-green-50 text-green-700 border-green-200',
    blue: 'bg-blue-50 text-blue-700 border-blue-200',
    purple: 'bg-purple-50 text-purple-700 border-purple-200',
    indigo: 'bg-indigo-50 text-indigo-700 border-indigo-200',
    yellow: 'bg-yellow-50 text-yellow-700 border-yellow-200',
  };
  return (
    <div className={`border-2 rounded-lg p-4 ${colors[color] ?? colors.blue}`}>
      <div className="text-xs font-medium opacity-75 mb-1">{label}</div>
      <div className="text-2xl font-bold">{value}</div>
    </div>
  );
}

function Center({ children, className = '' }: { children: React.ReactNode; className?: string }) {
  return (
    <div className={`min-h-screen bg-gray-50 flex items-center justify-center ${className}`}>
      <div className="text-gray-600">{children}</div>
    </div>
  );
}

function Empty({ children }: { children: React.ReactNode }) {
  return (
    <div className="bg-white rounded-lg shadow p-12 text-center text-gray-400">
      <p className="text-lg mb-2">No data yet</p>
      <p className="text-sm">{children}</p>
    </div>
  );
}
