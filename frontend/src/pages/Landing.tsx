import { Link } from 'react-router-dom';
import { Zap, Bell, Search, Tv } from 'lucide-react';

const features = [
  {
    icon: <Bell className="h-6 w-6" style={{ color: 'var(--kb-red)' }} />,
    title: 'Instant Telegram Alerts',
    body: 'Link your Telegram account and get notified the moment a deal lands.',
  },
  {
    icon: <Search className="h-6 w-6" style={{ color: 'var(--kb-red)' }} />,
    title: 'Keyword Watchlist',
    body: 'Watch for specific keywords like "iPad" or "Coffee Machine" — never miss a bargain.',
  },
];

export default function Landing() {
  return (
    <div className="flex min-h-screen flex-col" style={{ background: 'var(--kb-cream)' }}>
      {/* Nav */}
      <header
        className="flex items-center justify-between px-6 py-4"
        style={{ background: 'var(--kb-red)', color: '#fff' }}
      >
        <span className="serif flex items-center gap-2 text-2xl font-bold tracking-tight">
          <img src="/kramer-icon.jpg" alt="" className="h-9 w-9 rounded-full object-cover" />
          KramerBot - Aussie Deals
        </span>
        <div className="flex gap-3">
          <Link
            to="/login"
            className="rounded-lg px-4 py-2 text-sm font-semibold text-white hover:bg-white/20 transition"
          >
            Sign In
          </Link>
          <Link
            to="/signup"
            className="rounded-lg px-4 py-2 text-sm font-semibold transition"
            style={{ background: 'var(--kb-yellow)', color: 'var(--kb-ink)' }}
          >
            Get Started
          </Link>
        </div>
      </header>

      {/* Hero */}
      <main className="flex flex-1 flex-col items-center px-6 py-16 text-center">
        {/* Kramer photo */}
        <div className="mb-8 flex flex-col items-center">
          <div
            className="overflow-hidden rounded-2xl border-4 shadow-xl"
            style={{ borderColor: 'var(--kb-yellow)', maxWidth: 280 }}
          >
            <img
              src="/kramer.jpg"
              alt="Kramer"
              className="w-full object-cover"
              style={{ maxHeight: 340 }}
            />
          </div>
          <p
            className="serif mt-4 text-xl font-bold italic"
            style={{ color: 'var(--kb-red)' }}
          >
            "Looking for deals? Giddy up!"
          </p>
        </div>

        {/* Headline */}
        <h1
          className="serif mb-6 text-5xl font-extrabold tracking-tight sm:text-6xl"
          style={{ color: 'var(--kb-ink)' }}
        >
          The smartest deal hunter{' '}
          <span style={{ color: 'var(--kb-red)' }}>in Australia</span>
        </h1>
        <p className="mb-10 max-w-l text-md" style={{ color: '#444' }}>
          KramerBot monitors OzBargain and Amazon Australia 24/7, alerting you to hot deals via
          Telegram — filtered by upvotes and your personal keyword watchlist.
        </p>

        <div className="flex flex-wrap justify-center gap-4">
          <Link to="/signup" className="btn-red text-base">
            Create Free Account
          </Link>
        </div>
      </main>

      {/* Features */}
      <section className="mx-auto w-full max-w-4xl px-6 pb-24">
        <div className="grid gap-6 sm:grid-cols-2">
          {features.map((f) => (
            <div
              key={f.title}
              className="card"
            >
              <div
                className="mb-3 flex h-10 w-10 items-center justify-center rounded-xl"
                style={{ background: 'var(--kb-cream-dark)' }}
              >
                {f.icon}
              </div>
              <h3 className="serif mb-1 font-bold" style={{ color: 'var(--kb-ink)' }}>
                {f.title}
              </h3>
              <p className="text-sm text-slate-500">{f.body}</p>
            </div>
          ))}
        </div>
      </section>

      <footer
        className="border-t py-6 text-center text-sm"
        style={{ borderColor: 'var(--kb-yellow)', color: '#666' }}
      >
        KramerBot — <span style={{ color: 'var(--kb-red)' }}>Made with ❤️ in Adelaide, Australia 🤙</span>
      </footer>
    </div>
  );
}
