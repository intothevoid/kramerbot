import { useState } from 'react';
import { LogOut, Tag, Plus, X } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getOzbDeals,
  getAmazonDeals,
  getKeywords,
  addKeyword,
  removeKeyword,
  updatePreferences,
  getProfile,
} from '../api/user';
import { OzbDealCard, AmazonDealCard } from '../components/DealCard';
import { TelegramLinker } from '../components/TelegramLinker';
import { Spinner } from '../components/Spinner';
import type { WebUser } from '../types';

type Tab = 'ozb-good' | 'ozb-super' | 'amz-daily' | 'amz-weekly';

const tabs: { id: Tab; label: string }[] = [
  { id: 'ozb-good', label: '✅ OzBargain All' },
  { id: 'ozb-super', label: '⚡ OzBargain Top' },
  { id: 'amz-daily', label: '📅 Amazon Daily' },
  { id: 'amz-weekly', label: '📆 Amazon Weekly' },
];

interface Props {
  user: WebUser;
  onSignOut: () => void;
}

function ToggleRow({
  label,
  checked,
  onChange,
  disabled,
}: {
  label: string;
  checked: boolean;
  onChange: (v: boolean) => void;
  disabled?: boolean;
}) {
  return (
    <label className={`flex cursor-pointer items-center justify-between py-1.5 ${disabled ? 'opacity-50' : ''}`}>
      <span className="text-sm text-slate-700">{label}</span>
      <label className="toggle-switch">
        <input
          type="checkbox"
          checked={checked}
          onChange={(e) => !disabled && onChange(e.target.checked)}
          disabled={disabled}
        />
        <span className="toggle-slider" />
      </label>
    </label>
  );
}

export default function Dashboard({ user, onSignOut }: Props) {
  const qc = useQueryClient();
  const [tab, setTab] = useState<Tab>('ozb-good');
  const [newKeyword, setNewKeyword] = useState('');

  // Fresh profile to get prefs
  const profileQuery = useQuery({
    queryKey: ['profile'],
    queryFn: getProfile,
    initialData: user,
    staleTime: 30_000,
  });
  const profile = profileQuery.data ?? user;

  // Deal queries — all fetch on mount so tab switches are instant.
  // Responses come from in-memory scraper cache so the 4 requests are cheap.
  const ozbGood = useQuery({
    queryKey: ['deals', 'ozb', 'good'],
    queryFn: () => getOzbDeals(),
    staleTime: 2 * 60 * 1000,
  });
  const ozbSuper = useQuery({
    queryKey: ['deals', 'ozb', 'super'],
    queryFn: () => getOzbDeals('super'),
    staleTime: 2 * 60 * 1000,
  });
  const amzDaily = useQuery({
    queryKey: ['deals', 'amz', 'daily'],
    queryFn: () => getAmazonDeals('daily'),
    staleTime: 5 * 60 * 1000,
  });
  const amzWeekly = useQuery({
    queryKey: ['deals', 'amz', 'weekly'],
    queryFn: () => getAmazonDeals('weekly'),
    staleTime: 5 * 60 * 1000,
  });

  const keywordsQuery = useQuery({ queryKey: ['keywords'], queryFn: getKeywords });

  // Keyword mutations — update cache directly from server response
  const addKw = useMutation({
    mutationFn: addKeyword,
    onSuccess: (updated) => {
      setNewKeyword('');
      qc.setQueryData(['keywords'], updated);
    },
  });
  const removeKw = useMutation({
    mutationFn: removeKeyword,
    onSuccess: (updated) => qc.setQueryData(['keywords'], updated),
  });

  // Preference mutation
  const prefsMutation = useMutation({
    mutationFn: updatePreferences,
    onSuccess: (updated) => qc.setQueryData(['profile'], updated),
  });

  const handlePrefToggle = (key: 'ozb_good' | 'ozb_super' | 'amz_daily' | 'amz_weekly', value: boolean) => {
    prefsMutation.mutate({
      ozb_good: key === 'ozb_good' ? value : profile.ozb_good ?? false,
      ozb_super: key === 'ozb_super' ? value : profile.ozb_super ?? false,
      amz_daily: key === 'amz_daily' ? value : profile.amz_daily ?? false,
      amz_weekly: key === 'amz_weekly' ? value : profile.amz_weekly ?? false,
    });
  };

  const activeQuery = { 'ozb-good': ozbGood, 'ozb-super': ozbSuper, 'amz-daily': amzDaily, 'amz-weekly': amzWeekly }[tab];
  const isOzb = tab.startsWith('ozb');

  return (
    <div className="flex min-h-screen flex-col" style={{ background: 'var(--kb-cream)' }}>
      {/* Top nav */}
      <header
        className="flex items-center justify-between px-6 py-3 shadow-sm"
        style={{ background: 'var(--kb-red)' }}
      >
        <span className="serif flex items-center gap-2 text-lg font-bold text-white">
          <img src="/kramer-icon.jpg" alt="" className="h-8 w-8 rounded-full object-cover" />
          KramerBot - Aussie Deals
        </span>
        <div className="flex items-center gap-3">
          <span className="hidden text-sm text-red-100 sm:block">
            {profile.display_name || profile.email}
          </span>
          <button
            onClick={onSignOut}
            className="flex items-center gap-1 rounded-lg px-3 py-1.5 text-sm text-white hover:bg-red-700"
          >
            <LogOut className="h-4 w-4" /> Sign out
          </button>
        </div>
      </header>

      <div className="mx-auto flex w-full max-w-6xl flex-1 gap-6 px-4 py-6 lg:px-6">
        {/* Sidebar */}
        <aside className="hidden w-72 shrink-0 space-y-5 lg:block">

          {/* Telegram */}
          <div className="card">
            <h3
              className="serif mb-3 text-sm font-bold uppercase tracking-wide"
              style={{ color: 'var(--kb-red)' }}
            >
              Telegram
            </h3>
            <TelegramLinker />
          </div>

          {/* Subscriptions */}
          <div className="card">
            <h3
              className="serif mb-1 text-sm font-bold uppercase tracking-wide"
              style={{ color: 'var(--kb-red)' }}
            >
              Subscriptions
            </h3>
            <p className="mb-3 text-xs text-slate-500">
              Choose which deals Telegram sends you.
            </p>
            {prefsMutation.isPending && (
              <div className="mb-2 flex items-center gap-1 text-xs text-slate-400">
                <Spinner size="sm" /> Saving…
              </div>
            )}
            <div className="divide-y divide-slate-100">
              <ToggleRow
                label="✅ OzBargain All deals"
                checked={profile.ozb_good ?? false}
                onChange={(v) => handlePrefToggle('ozb_good', v)}
                disabled={prefsMutation.isPending}
              />
              <ToggleRow
                label="⚡ OzBargain Top (25+ votes)"
                checked={profile.ozb_super ?? false}
                onChange={(v) => handlePrefToggle('ozb_super', v)}
                disabled={prefsMutation.isPending}
              />
              <ToggleRow
                label="📅 Amazon Daily"
                checked={profile.amz_daily ?? false}
                onChange={(v) => handlePrefToggle('amz_daily', v)}
                disabled={prefsMutation.isPending}
              />
              <ToggleRow
                label="📆 Amazon Weekly"
                checked={profile.amz_weekly ?? false}
                onChange={(v) => handlePrefToggle('amz_weekly', v)}
                disabled={prefsMutation.isPending}
              />
            </div>
          </div>

          {/* Keywords */}
          <div className="card">
            <h3
              className="serif mb-3 flex items-center gap-1 text-sm font-bold uppercase tracking-wide"
              style={{ color: 'var(--kb-red)' }}
            >
              <Tag className="h-3.5 w-3.5" /> Keywords
            </h3>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                if (newKeyword.trim()) addKw.mutate(newKeyword.trim());
              }}
              className="mb-3 flex gap-2"
            >
              <input
                value={newKeyword}
                onChange={(e) => setNewKeyword(e.target.value)}
                placeholder="e.g. iPad"
                className="input-field flex-1"
              />
              <button
                type="submit"
                disabled={addKw.isPending}
                className="btn-red px-3"
                style={{ borderRadius: '0.5rem' }}
              >
                <Plus className="h-4 w-4" />
              </button>
            </form>
            <div className="flex flex-wrap gap-1.5">
              {keywordsQuery.data?.map((kw) => (
                <span
                  key={kw}
                  className="flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium"
                  style={{ background: 'var(--kb-yellow)', color: 'var(--kb-ink)' }}
                >
                  {kw}
                  <button
                    onClick={() => removeKw.mutate(kw)}
                    className="hover:opacity-70"
                    disabled={removeKw.isPending}
                  >
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
                  tab === t.id ? 'text-white shadow' : 'text-slate-600 hover:bg-slate-100'
                }`}
                style={tab === t.id ? { background: 'var(--kb-red)' } : {}}
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
              <div key={tab} className="grid gap-3 sm:grid-cols-2">
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
