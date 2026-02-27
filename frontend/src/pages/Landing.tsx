import { Link } from 'react-router-dom';
import { Zap, Bell, Search, Tv } from 'lucide-react';

const features = [
  {
    icon: <Zap className="h-6 w-6 text-indigo-500" />,
    title: 'Real-time Deals',
    body: 'Scraped every 5 minutes from OzBargain and Amazon Australia.',
  },
  {
    icon: <Bell className="h-6 w-6 text-indigo-500" />,
    title: 'Instant Telegram Alerts',
    body: 'Link your Telegram account and get notified the moment a deal lands.',
  },
  {
    icon: <Search className="h-6 w-6 text-indigo-500" />,
    title: 'Keyword Watchlist',
    body: 'Watch for specific keywords like "iPad" or "Coffee Machine" and never miss a bargain.',
  },
  {
    icon: <Tv className="h-6 w-6 text-indigo-500" />,
    title: 'Android TV Notifications',
    body: 'Optionally push deals straight to your TV via Pipup.',
  },
];

export default function Landing() {
  return (
    <div className="flex min-h-screen flex-col bg-gradient-to-br from-slate-50 to-indigo-50">
      {/* Nav */}
      <header className="flex items-center justify-between border-b border-slate-200 bg-white/80 px-6 py-4 backdrop-blur">
        <span className="text-xl font-bold text-indigo-700">KramerBot 🛒</span>
        <div className="flex gap-3">
          <Link to="/login" className="rounded-lg px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100">
            Sign In
          </Link>
          <Link
            to="/signup"
            className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
          >
            Get Started
          </Link>
        </div>
      </header>

      {/* Hero */}
      <main className="flex flex-1 flex-col items-center justify-center px-6 py-20 text-center">
        <span className="mb-4 rounded-full border border-indigo-200 bg-indigo-50 px-4 py-1 text-sm font-medium text-indigo-700">
          Let Kramer watch deals so you don't have to
        </span>
        <h1 className="mb-6 text-5xl font-extrabold tracking-tight text-slate-900 sm:text-6xl">
          The smartest deal hunter <br className="hidden sm:block" />
          <span className="text-indigo-600">in Australia</span>
        </h1>
        <p className="mb-10 max-w-xl text-lg text-slate-600">
          KramerBot monitors OzBargain and Amazon Australia 24/7, alerting you to hot deals based on upvotes
          and your personal keyword watchlist — delivered instantly via Telegram.
        </p>
        <div className="flex flex-wrap justify-center gap-4">
          <Link
            to="/signup"
            className="rounded-xl bg-indigo-600 px-8 py-3.5 text-base font-semibold text-white shadow hover:bg-indigo-700"
          >
            Create Free Account
          </Link>
          <a
            href="https://t.me/kramerbot"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-xl border border-slate-300 bg-white px-8 py-3.5 text-base font-semibold text-slate-700 hover:border-slate-400"
          >
            Open in Telegram
          </a>
        </div>
      </main>

      {/* Features */}
      <section className="mx-auto w-full max-w-4xl px-6 pb-24">
        <div className="grid gap-6 sm:grid-cols-2">
          {features.map((f) => (
            <div key={f.title} className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-50">
                {f.icon}
              </div>
              <h3 className="mb-1 font-semibold text-slate-900">{f.title}</h3>
              <p className="text-sm text-slate-500">{f.body}</p>
            </div>
          ))}
        </div>
      </section>

      <footer className="border-t border-slate-200 py-6 text-center text-sm text-slate-400">
        KramerBot — Giddy up! 🤙
      </footer>
    </div>
  );
}
