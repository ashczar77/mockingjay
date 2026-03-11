'use client';

import { useEffect, useState } from 'react';

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
}

export default function Home() {
  const [metrics, setMetrics] = useState<ConversationMetrics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/metrics')
      .then(res => res.json())
      .then(data => {
        setMetrics(data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to fetch metrics:', err);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (!metrics) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-red-600">Failed to load metrics</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-2">🐦 MockingJay</h1>
        <p className="text-gray-600 mb-8">Voice AI Testing Platform - Conversation Intelligence</p>

        {/* Conversation Intelligence */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-2xl font-semibold mb-4">💬 Conversation Intelligence</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <MetricCard
              label="Success Rate"
              value={`${metrics.success_rate.toFixed(1)}%`}
              color="green"
            />
            <MetricCard
              label="Intent Accuracy"
              value={`${metrics.intent_accuracy.toFixed(1)}%`}
              color="blue"
            />
            <MetricCard
              label="Avg Steps"
              value={metrics.avg_steps_completed.toFixed(1)}
              color="purple"
            />
            <MetricCard
              label="Multi-turn"
              value={metrics.multi_turn_count.toString()}
              color="indigo"
            />
          </div>
        </div>

        {/* Multi-turn Dialogue */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-2xl font-semibold mb-4">🔄 Multi-turn Dialogue</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            <MetricCard
              label="Context Retention"
              value={`${metrics.context_retention.toFixed(1)}%`}
              color="green"
            />
            <MetricCard
              label="Coherence Score"
              value={`${metrics.coherence_score.toFixed(1)}%`}
              color="blue"
            />
            <MetricCard
              label="Multi-turn Count"
              value={metrics.multi_turn_count.toString()}
              color="purple"
            />
          </div>
        </div>

        {/* Response Quality */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-4">✨ Response Quality</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <MetricCard
              label="Completeness"
              value={`${metrics.completeness_score.toFixed(1)}%`}
              color="green"
            />
            <MetricCard
              label="Sentiment"
              value={`${metrics.sentiment_score.toFixed(1)}%`}
              color="yellow"
            />
            <MetricCard
              label="Confidence"
              value={`${metrics.confidence_score.toFixed(1)}%`}
              color="blue"
            />
            <MetricCard
              label="Avg Length"
              value={`${metrics.avg_response_length.toFixed(0)} chars`}
              color="purple"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

interface MetricCardProps {
  label: string;
  value: string;
  color: 'green' | 'blue' | 'purple' | 'indigo' | 'yellow';
}

function MetricCard({ label, value, color }: MetricCardProps) {
  const colorClasses = {
    green: 'bg-green-50 text-green-700 border-green-200',
    blue: 'bg-blue-50 text-blue-700 border-blue-200',
    purple: 'bg-purple-50 text-purple-700 border-purple-200',
    indigo: 'bg-indigo-50 text-indigo-700 border-indigo-200',
    yellow: 'bg-yellow-50 text-yellow-700 border-yellow-200',
  };

  return (
    <div className={`border-2 rounded-lg p-4 ${colorClasses[color]}`}>
      <div className="text-sm font-medium opacity-75 mb-1">{label}</div>
      <div className="text-2xl font-bold">{value}</div>
    </div>
  );
}
