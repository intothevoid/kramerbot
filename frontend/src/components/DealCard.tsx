import { ExternalLink, ThumbsUp, Clock } from 'lucide-react';
import type { OzbDeal, AmazonDeal } from '../types';

function OzbBadge({ type }: { type: number }) {
  // dealtype: 3 = OZB_GOOD, 2 = OZB_SUPER
  if (type === 2)
    return <span className="rounded-full bg-yellow-100 px-2 py-0.5 text-xs font-semibold text-yellow-800">⚡ Super</span>;
  return <span className="rounded-full bg-green-100 px-2 py-0.5 text-xs font-semibold text-green-800">✅ Good</span>;
}

export function OzbDealCard({ deal }: { deal: OzbDeal }) {
  return (
    <a
      href={deal.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group flex flex-col gap-2 rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:border-indigo-300 hover:shadow-md"
    >
      <div className="flex items-start justify-between gap-2">
        <p className="line-clamp-2 font-semibold text-slate-800 group-hover:text-indigo-700">
          {deal.title}
        </p>
        <ExternalLink className="mt-0.5 h-4 w-4 shrink-0 text-slate-400 group-hover:text-indigo-500" />
      </div>
      <div className="flex flex-wrap items-center gap-2 text-sm text-slate-500">
        <OzbBadge type={deal.dealtype} />
        <span className="flex items-center gap-1">
          <ThumbsUp className="h-3.5 w-3.5" />
          {deal.upvotes}
        </span>
        <span className="flex items-center gap-1">
          <Clock className="h-3.5 w-3.5" />
          {deal.dealage || deal.time}
        </span>
      </div>
    </a>
  );
}

export function AmazonDealCard({ deal }: { deal: AmazonDeal }) {
  const label = deal.dealtype === 4 ? '📅 Daily' : '📆 Weekly';
  return (
    <a
      href={deal.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group flex items-start gap-3 rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:border-indigo-300 hover:shadow-md"
    >
      {deal.image && (
        <img
          src={deal.image}
          alt=""
          className="h-16 w-16 shrink-0 rounded-lg object-contain"
          onError={(e) => (e.currentTarget.style.display = 'none')}
        />
      )}
      <div className="flex flex-col gap-1 overflow-hidden">
        <p className="line-clamp-2 font-semibold text-slate-800 group-hover:text-indigo-700">
          {deal.title}
        </p>
        <div className="flex items-center gap-2 text-sm text-slate-500">
          <span className="rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-800">
            {label}
          </span>
          <ExternalLink className="h-3.5 w-3.5" />
        </div>
      </div>
    </a>
  );
}
