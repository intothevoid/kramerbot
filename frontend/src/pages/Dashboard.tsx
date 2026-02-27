import { useState } from 'react';
import { LogOut, Tag, Plus, X } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getOzbDeals, getAmazonDeals, getKeywords, addKeyword, removeKeyword } from '../api/user';
import { OzbDealCard, AmazonDealCard } from '../components/DealCard';
import { TelegramLinker } from '../components/TelegramLinker';
import { Spinner } from '../components/Spinner';
import type { WebUser } from '../types';

type Tab = 'ozb-good' | 'ozb-super' | 'amz-daily' | 'amz-weekly';

const tabs: { id: Tab; label: string }[] = [
  { id: 'ozb-good', label: '✅ OzBargain Good' },
  { id: 'ozb-super', label: '⚡ OzBargain Super' },
  { id: 'amz-daily', label: '📅 Amazon Daily' },
  { id: 'amz-weekly', label: '📆 Amazon Weekly' },
];

interface Props {
  user: WebUser;
  onSignOut: () => void;
}

export default function Dashboard({ user, onSignOut }: Props) {
  const qc = useQueryClient();
  const [tab, setTab] = useState<Tab>('ozb-good');
  const [newKeyword, setNewKeyword] = useState('');

  // Deal queries
  const ozbGood = useQuery({
    queryKey: ['deals', 'ozb', 'good'],
    queryFn: () => getOzbDeals('good'),
    enabled: tab === 'ozb-good',
    staleTime: 2 * 60 * 1000,
  });
  const ozbSuper = useQuery({
    queryKey: ['deals', 'ozb', 'super'],
    queryFn: () => getOzbDeals('super'),
    enabled: tab === 'ozb-super',
    staleTime: 2 * 60 * 1000,
  });
  const amzDaily = useQuery({
    queryKey: ['deals', 'amz', 'daily'],
    queryFn: () => getAmazonDeals('daily'),
    enabled: tab === 'amz-daily',
    staleTime: 5 * 60 * 1000,
  });
  const amzWeekly = useQuery({
    queryKey: ['deals', 'amz', 'weekly'],
    queryFn: () => getAmazonDeals('weekly'),
    enabled: tab === 'amz-weekly',
    staleTime: 5 * 60 * 1000,
  });

  const keywordsQuery = useQuery({ queryKey: ['keywords'], queryFn: getKeywords });

  const addKw = useMutation({
    mutationFn: addKeyword,
    onSuccess: () => {
      setNewKeyword('');
      qc.invalidateQueries({ queryKey: ['keywords'] });
    },
  });
  const removeKw = useMutation({
    mutationFn: removeKeyword,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['keywords'] }),
  });

  const activeQuery = { 'ozb-good': ozbGood, 'ozb-super': ozbSuper, 'amz-daily': amzDaily, 'amz-weekly': amzWeekly }[tab];
  const isOzb = tab.startsWith('ozb');

  return (
    <div className="flex min-h-screen flex-col bg-slate-50">
      {/* Top nav */}
      <header className="flex items-center justify-between border-b border-slate-200 bg-white px-6 py-3">
        <span className="text-lg font-bold text-indigo-700">KramerBot 🛒</span>
        <div className="flex items-center gap-3">
          <span className="hidden text-sm text-slate-500 sm:block">
            {user.display_name || user.email}
          </span>
          <button
            onClick={onSignOut}
            className="flex items-center gap-1 rounded-lg px-3 py-1.5 text-sm text-slate-600 hover:bg-slate-100"
          >
            <LogOut className="h-4 w-4" /> Sign out
          </button>
        </div>
      </header>

      <div className="mx-auto flex w-full max-w-6xl flex-1 gap-6 px-4 py-6 lg:px-6">
        {/* Sidebar */}
        <aside className="hidden w-64 shrink-0 space-y-5 lg:block">
          {/* Telegram */}
          <div>
            <h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">Telegram</h3>
            <TelegramLinker />
          </div>

          {/* Keywords */}
          <div>
            <h3 className="mb-2 flex items-center gap-1 text-xs font-semibold uppercase tracking-wide text-slate-500">
              <Tag className="h-3.5 w-3.5" /> Keywords
            </h3>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                if (newKeyword.trim()) addKw.mutate(newKeyword.trim());
              }}
              className="mb-2 flex gap-2"
            >
              <input
                value={newKeyword}
                onChange={(e) => setNewKeyword(e.target.value)}
                placeholder="e.g. iPad"
                className="flex-1 rounded-lg border border-slate-300 px-3 py-1.5 text-sm outline-none focus:border-indigo-400"
              />
              <button
                type="submit"
                disabled={addKw.isPending}
                className="rounded-lg bg-indigo-600 px-2.5 text-white hover:bg-indigo-700 disabled:opacity-60"
              >
                <Plus className="h-4 w-4" />
              </button>
            </form>
            <div className="flex flex-wrap gap-1.5">
              {keywordsQuery.data?.map((kw) => (
                <span key={kw} className="flex items-center gap-1 rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-700">
                  {kw}
                  <button onClick={() => removeKw.mutate(kw)} className="hover:text-red-600">
                    <X className="h-3 w-3" />
                  </button>
                </span>
              ))}
              {keywordsQuery.data?.length === 0 && (
                <span className="text-xs text-slate-400">No keywords yet</span>
              )}
            </div>
          </div>
        </aside>

        {/* Main deal feed */}
        <main className="flex-1">
          {/* Tabs */}
          <div className="mb-4 flex gap-1 overflow-x-auto rounded-xl bg-white p-1 shadow-sm">
            {tabs.map((t) => (
              <button
                key={t.id}
                onClick={() => setTab(t.id)}
                className={`whitespace-nowrap rounded-lg px-4 py-2 text-sm font-medium transition ${
                  tab === t.id
                    ? 'bg-indigo-600 text-white shadow'
                    : 'text-slate-600 hover:bg-slate-100'
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>

          {/* Deal list */}
          {activeQuery.isLoading && (
            <div className="flex justify-center py-16">
              <Spinner size="lg" />
            </div>
          )}
          {activeQuery.isError && (
            <div className="rounded-xl bg-red-50 p-6 text-center text-red-700">
              Failed to load deals. The scrapers may still be warming up.
            </div>
          )}
          {activeQuery.data && (
            <>
              <p className="mb-3 text-xs text-slate-400">{activeQuery.data.total} deals</p>
              <div className="grid gap-3 sm:grid-cols-2">
                {isOzb
                  ? activeQuery.data.deals.map((d: unknown) => {
                      const deal = d as import('../types').OzbDeal;
                      return <OzbDealCard key={deal.id} deal={deal} />;
                    })
                  : activeQuery.data.deals.map((d: unknown) => {
                      const deal = d as import('../types').AmazonDeal;
                      return <AmazonDealCard key={deal.id} deal={deal} />;
                    })}
              </div>
              {activeQuery.data.deals.length === 0 && (
                <div className="rounded-xl bg-white p-10 text-center text-slate-500 shadow-sm">
                  No deals yet — check back in a few minutes while the scrapers warm up.
                </div>
              )}
            </>
          )}
        </main>
      </div>
    </div>
  );
}
