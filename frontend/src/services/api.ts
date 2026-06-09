const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export interface ETF {
  id: number;
  symbol: string;
  name: string;
  description: string;
  expense_ratio: number;
  created_at: string;
}

export interface Price {
  id: number;
  etf_id: number;
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  created_at: string;
}

export interface Platform {
  id: number;
  name: string;
  description: string;
  fees: string;
  pros: string[];
  cons: string[];
  rating: number;
  website: string;
}

export interface News {
  id: number;
  title: string;
  url: string;
  published_at: string;
  source: string;
  summary: string;
}

export interface PriceAlert {
  id: number;
  etf_id: number;
  email: string;
  threshold: number;
  direction: string;
  created_at: string;
  triggered: boolean;
}

export const api = {
  async getETFs(): Promise<ETF[]> {
    const response = await fetch(`${API_URL}/etfs`);
    if (!response.ok) throw new Error('Failed to fetch ETFs');
    return response.json();
  },

  async getETFBySymbol(symbol: string): Promise<ETF> {
    const response = await fetch(`${API_URL}/etfs/${symbol}`);
    if (!response.ok) throw new Error('Failed to fetch ETF');
    return response.json();
  },

  async getETFPrices(symbol: string, limit: number = 365): Promise<Price[]> {
    const response = await fetch(`${API_URL}/etfs/${symbol}/prices?limit=${limit}`);
    if (!response.ok) throw new Error('Failed to fetch ETF prices');
    return response.json();
  },

  async getTopPerformers(period: string = '1y', limit: number = 10): Promise<ETF[]> {
    const response = await fetch(`${API_URL}/etfs/top?period=${period}&limit=${limit}`);
    if (!response.ok) throw new Error('Failed to fetch top performers');
    return response.json();
  },

  async getPlatforms(): Promise<Platform[]> {
    const response = await fetch(`${API_URL}/platforms`);
    if (!response.ok) throw new Error('Failed to fetch platforms');
    return response.json();
  },

  async getNews(limit: number = 20): Promise<News[]> {
    const response = await fetch(`${API_URL}/news?limit=${limit}`);
    if (!response.ok) throw new Error('Failed to fetch news');
    return response.json();
  },

  async createAlert(alert: Omit<PriceAlert, 'id' | 'created_at' | 'triggered'>): Promise<PriceAlert> {
    const response = await fetch(`${API_URL}/alerts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(alert),
    });
    if (!response.ok) throw new Error('Failed to create alert');
    return response.json();
  },
};
