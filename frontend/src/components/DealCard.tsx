import { ExternalLink, ThumbsUp, Clock } from 'lucide-react';
import type { OzbDeal, AmazonDeal } from '../types';

function OzbBadge({ type }: { type: number }) {
  // dealtype: 2 = OZB_SUPER (top deal: 25+ votes in 1h), 1 = OZB_REG
  if (type === 2)
    return (
      <span
        className="rounded-full px-2 py-0.5 text-xs font-semibold"
        style={{ background: 'var(--kb-yellow)', color: 'var(--kb-ink)' }}
      >
        ⚡ Top
      </span>
    );
  return (
    <span
      className="rounded-full px-2 py-0.5 text-xs font-semibold"
      style={{ background: '#f1f5f9', color: '#475569' }}
    >
      Deal
    </span>
  );
}

export function OzbDealCard({ deal }: { deal: OzbDeal }) {
  return (
    <a
      href={deal.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group card flex flex-col gap-2 transition hover:shadow-md"
      style={{ borderColor: 'transparent' }}
      onMouseEnter={(e) => (e.currentTarget.style.borderColor = 'var(--kb-yellow)')}
      onMouseLeave={(e) => (e.currentTarget.style.borderColor = 'transparent')}
    >
      <div className="flex items-start justify-between gap-2">
        <p className="line-clamp-2 font-semibold text-slate-800 group-hover:text-red-700">
          {deal.title}
        </p>
        <ExternalLink className="mt-0.5 h-4 w-4 shrink-0 text-slate-400 group-hover:text-red-500" />
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
  const isDaily = deal.dealtype === 4;
  return (
    <a
      href={deal.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group card flex items-start gap-3 transition hover:shadow-md"
      style={{ borderColor: 'transparent' }}
      onMouseEnter={(e) => (e.currentTarget.style.borderColor = 'var(--kb-yellow)')}
      onMouseLeave={(e) => (e.currentTarget.style.borderColor = 'transparent')}
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
        <p className="line-clamp-2 font-semibold text-slate-800 group-hover:text-red-700">
          {deal.title}
        </p>
        <div className="flex items-center gap-2 text-sm text-slate-500">
          <span
            className="rounded-full px-2 py-0.5 text-xs font-semibold"
            style={{ background: 'var(--kb-yellow)', color: 'var(--kb-ink)' }}
          >
            {isDaily ? '📅 Daily' : '📆 Weekly'}
          </span>
          <ExternalLink className="h-3.5 w-3.5" />
        </div>
      </div>
    </a>
  );
}
