import { useState, useEffect } from 'react';
import { ETFList } from './components/ETFList';
import { PlatformReviews } from './components/PlatformReviews';
import { NewsFeed } from './components/NewsFeed';
import { PriceAlert } from './components/PriceAlert';
import { ETFDetail } from './components/ETFDetail';
import { api, type ETF } from './services/api';
import { TrendingUp, LayoutDashboard, Newspaper, Bell } from 'lucide-react';

function App() {
  const [activeTab, setActiveTab] = useState<'etfs' | 'platforms' | 'news' | 'alerts'>('etfs');
  const [selectedETF, setSelectedETF] = useState<string | null>(null);
  const [etfs, setEtfs] = useState<ETF[]>([]);

  useEffect(() => {
    loadETFs();
  }, []);

  const loadETFs = async () => {
    try {
      const data = await api.getETFs();
      setEtfs(data);
    } catch (err) {
      console.error('Failed to load ETFs:', err);
    }
  };

  const tabs = [
    { id: 'etfs' as const, label: 'ETFs', icon: TrendingUp },
    { id: 'platforms' as const, label: 'Platforms', icon: LayoutDashboard },
    { id: 'news' as const, label: 'News', icon: Newspaper },
    { id: 'alerts' as const, label: 'Alerts', icon: Bell },
  ];

  if (selectedETF) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <ETFDetail 
            symbol={selectedETF} 
            onClose={() => setSelectedETF(null)} 
          />
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            Philippine ETF Recommendations
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            Real-time ETF analysis, platform reviews, and investment insights
          </p>
        </div>
      </header>

      <nav className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4">
          <div className="flex gap-1 overflow-x-auto py-2">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition whitespace-nowrap ${
                    activeTab === tab.id
                      ? 'bg-blue-600 text-white'
                      : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
                  }`}
                >
                  <Icon className="w-4 h-4" />
                  {tab.label}
                </button>
              );
            })}
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 py-8">
        {activeTab === 'etfs' && <ETFList />}
        {activeTab === 'platforms' && <PlatformReviews />}
        {activeTab === 'news' && <NewsFeed />}
        {activeTab === 'alerts' && <PriceAlert etfs={etfs} />}
      </main>

      <footer className="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 mt-12">
        <div className="max-w-7xl mx-auto px-4 py-6 text-center text-gray-600 dark:text-gray-400">
          <p>© 2024 Philippine ETF Recommendations. Data provided for informational purposes only.</p>
        </div>
      </footer>
    </div>
  );
}

export default App;
