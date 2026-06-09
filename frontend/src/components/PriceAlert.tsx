import { useState } from 'react';
import { api, type ETF, type PriceAlert } from '../services/api';
import { Bell, AlertCircle, CheckCircle } from 'lucide-react';

export function PriceAlert({ etfs }: { etfs: ETF[] }) {
  const [selectedETF, setSelectedETF] = useState<string>('');
  const [email, setEmail] = useState('');
  const [threshold, setThreshold] = useState('');
  const [direction, setDirection] = useState<'above' | 'below'>('above');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!selectedETF || !email || !threshold) {
      setError('Please fill in all fields');
      return;
    }

    const etf = etfs.find(e => e.symbol === selectedETF);
    if (!etf) {
      setError('Invalid ETF selected');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      await api.createAlert({
        etf_id: etf.id,
        email,
        threshold: parseFloat(threshold),
        direction,
      });

      setSuccess(true);
      // Reset form
      setSelectedETF('');
      setEmail('');
      setThreshold('');
      
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      setError('Failed to create alert. Please try again.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700">
      <div className="flex items-center gap-3 mb-6">
        <Bell className="w-6 h-6 text-blue-600" />
        <h2 className="text-xl font-bold text-gray-900 dark:text-white">
          Set Price Alert
        </h2>
      </div>

      {success && (
        <div className="mb-4 p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg flex items-center gap-3">
          <CheckCircle className="w-5 h-5 text-green-600" />
          <p className="text-sm text-green-800 dark:text-green-200">
            Alert created successfully! You'll receive an email when the price threshold is reached.
          </p>
        </div>
      )}

      {error && (
        <div className="mb-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-center gap-3">
          <AlertCircle className="w-5 h-5 text-red-600" />
          <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Select ETF
          </label>
          <select
            value={selectedETF}
            onChange={(e) => setSelectedETF(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          >
            <option value="">Choose an ETF...</option>
            {etfs.map((etf) => (
              <option key={etf.id} value={etf.symbol}>
                {etf.symbol} - {etf.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Email Address
          </label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="your@email.com"
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Price Threshold (PHP)
          </label>
          <input
            type="number"
            step="0.01"
            value={threshold}
            onChange={(e) => setThreshold(e.target.value)}
            placeholder="100.00"
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Alert When Price Is
          </label>
          <div className="flex gap-4">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                value="above"
                checked={direction === 'above'}
                onChange={(e) => setDirection(e.target.value as 'above' | 'below')}
                className="w-4 h-4 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-gray-700 dark:text-gray-300">Above</span>
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                value="below"
                checked={direction === 'below'}
                onChange={(e) => setDirection(e.target.value as 'above' | 'below')}
                className="w-4 h-4 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-gray-700 dark:text-gray-300">Below</span>
            </label>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed font-medium"
        >
          {loading ? 'Creating Alert...' : 'Create Alert'}
        </button>
      </form>

      <p className="mt-4 text-xs text-gray-500 dark:text-gray-400">
        Alerts are checked daily. You'll receive an email notification when the price threshold is reached.
      </p>
    </div>
  );
}
