import { useEffect, useState } from 'react';
import { api, type ETF } from '../services/api';
import { TrendingUp, DollarSign } from 'lucide-react';

export function ETFList() {
  const [etfs, setETFs] = useState<ETF[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadETFs();
  }, []);

  const loadETFs = async () => {
    try {
      setLoading(true);
      const data = await api.getETFs();
      setETFs(data);
    } catch (err) {
      setError('Failed to load ETFs');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-gray-500">Loading ETFs...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
          Philippine ETFs
        </h2>
        <button
          onClick={loadETFs}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
        >
          Refresh
        </button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {etfs.map((etf) => (
          <div
            key={etf.id}
            className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition"
          >
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-xl font-bold text-gray-900 dark:text-white">
                  {etf.symbol}
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                  {etf.name}
                </p>
              </div>
              <div className="flex items-center gap-2 text-green-600">
                <TrendingUp className="w-5 h-5" />
              </div>
            </div>

            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
              {etf.description}
            </p>

            <div className="flex items-center justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                <DollarSign className="w-4 h-4" />
                <span>Expense Ratio: {(etf.expense_ratio * 100).toFixed(2)}%</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
