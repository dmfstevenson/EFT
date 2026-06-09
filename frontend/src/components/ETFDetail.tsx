import { useEffect, useState } from 'react';
import { api, type ETF, type Price } from '../services/api';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { TrendingUp, TrendingDown, DollarSign, Calendar } from 'lucide-react';

interface ETFDetailProps {
  symbol: string;
  onClose: () => void;
}

export function ETFDetail({ symbol, onClose }: ETFDetailProps) {
  const [etf, setEtf] = useState<ETF | null>(null);
  const [prices, setPrices] = useState<Price[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadETFData();
  }, [symbol]);

  const loadETFData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [etfData, pricesData] = await Promise.all([
        api.getETFBySymbol(symbol),
        api.getETFPrices(symbol, 365),
      ]);
      
      setEtf(etfData);
      setPrices(pricesData.reverse()); // Show oldest to newest for chart
    } catch (err) {
      setError('Failed to load ETF data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const chartData = prices.map(p => ({
    date: new Date(p.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    price: p.close,
  }));

  const calculatePerformance = () => {
    if (prices.length < 2) return null;
    const first = prices[0].close;
    const last = prices[prices.length - 1].close;
    const change = ((last - first) / first) * 100;
    return { change, first, last };
  };

  const performance = calculatePerformance();

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-gray-500">Loading ETF details...</div>
      </div>
    );
  }

  if (error || !etf) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-red-500">{error || 'ETF not found'}</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold text-gray-900 dark:text-white">
            {etf.symbol}
          </h2>
          <p className="text-lg text-gray-600 dark:text-gray-400 mt-1">
            {etf.name}
          </p>
        </div>
        <button
          onClick={onClose}
          className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
        >
          Close
        </button>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700">
        <p className="text-gray-600 dark:text-gray-400 mb-6">{etf.description}</p>

        {performance && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                <Calendar className="w-4 h-4" />
                <span>Price Range</span>
              </div>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                ₱{performance.first.toFixed(2)} - ₱{performance.last.toFixed(2)}
              </p>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                {performance.change >= 0 ? (
                  <TrendingUp className="w-4 h-4 text-green-600" />
                ) : (
                  <TrendingDown className="w-4 h-4 text-red-600" />
                )}
                <span>1-Year Performance</span>
              </div>
              <p className={`text-2xl font-bold ${performance.change >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {performance.change >= 0 ? '+' : ''}{performance.change.toFixed(2)}%
              </p>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
                <DollarSign className="w-4 h-4" />
                <span>Expense Ratio</span>
              </div>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {(etf.expense_ratio * 100).toFixed(2)}%
              </p>
            </div>
          </div>
        )}

        {chartData.length > 0 && (
          <div className="h-80">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Price History (Last 365 Days)
            </h3>
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="date" 
                  stroke="#6B7280"
                  fontSize={12}
                />
                <YAxis 
                  stroke="#6B7280"
                  fontSize={12}
                  tickFormatter={(value) => `₱${value.toFixed(0)}`}
                />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#1F2937', 
                    border: '1px solid #374151',
                    borderRadius: '8px',
                  }}
                  itemStyle={{ color: '#F3F4F6' }}
                  labelStyle={{ color: '#9CA3AF' }}
                />
                <Line 
                  type="monotone" 
                  dataKey="price" 
                  stroke="#3B82F6" 
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        )}
      </div>
    </div>
  );
}
