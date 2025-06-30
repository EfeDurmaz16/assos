/* YouTube Heritage */
--youtube-red: #FF0000        /* YouTube Primary */
--youtube-dark: #CC0000       /* YouTube Hover */
--youtube-light: #FF4444      /* YouTube Light */

/* Off-White Variants */
--paper-white: #FAFAFA        /* Background */
--pearl-white: #F8F8F8        /* Cards */
--cream-white: #FFF8F0        /* Warm accent */

/* AI Purple Gradient */
--ai-purple-deep: #6B46C1     /* Deep AI */
--ai-purple-main: #8B5CF6     /* Main AI */
--ai-purple-light: #A78BFA    /* Light AI */
--ai-purple-glow: #C4B5FD     /* Glow effect */

/* Automation Green */
--auto-green-dark: #059669    /* Success dark */
--auto-green-main: #10B981    /* Success main */
--auto-green-light: #34D399   /* Success light */
--auto-green-glow: #6EE7B7    /* Processing */

/* Neutral Grays */
--gray-900: #111827           /* Text primary */
--gray-700: #374151           /* Text secondary */
--gray-500: #6B7280           /* Text muted */
--gray-300: #D1D5DB           /* Borders */
--gray-100: #F3F4F6           /* Backgrounds */

/* Status Colors */
--status-processing: #8B5CF6   /* AI Purple - Processing */
--status-active: #10B981       /* Green - Active */
--status-scheduled: #F59E0B    /* Amber - Scheduled */
--status-error: #EF4444        /* Red - Error */
--status-paused: #6B7280       /* Gray - Paused */

/* Interactive States */
--hover-red: #DC2626          /* Darker YouTube red */
--focus-purple: #7C3AED       /* Purple focus ring */
--active-green: #047857       /* Pressed state */

🎯 UI COMPONENT EXAMPLES
Primary Button (YouTube Style)
css.btn-primary {
  background: var(--youtube-red);
  color: var(--paper-white);
  border: none;
  padding: 12px 24px;
  border-radius: 2px;
  transition: all 0.2s;
}

.btn-primary:hover {
  background: var(--youtube-dark);
  transform: translateY(-1px);
}
AI Processing Card
css.ai-card {
  background: var(--pearl-white);
  border: 1px solid var(--gray-300);
  border-radius: 8px;
  position: relative;
  overflow: hidden;
}

.ai-card.processing::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 2px;
  background: linear-gradient(
    90deg, 
    transparent,
    var(--ai-purple-main),
    var(--ai-purple-light),
    transparent
  );
  animation: ai-scan 2s linear infinite;
}
Status Indicators
css.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.processing {
  background: var(--ai-purple-light);
  color: var(--ai-purple-deep);
}

.status-badge.active {
  background: var(--auto-green-light);
  color: var(--auto-green-dark);
}
Gradient Overlays
css/* AI Processing Gradient */
.ai-gradient {
  background: linear-gradient(
    135deg,
    var(--ai-purple-deep) 0%,
    var(--ai-purple-main) 50%,
    var(--ai-purple-light) 100%
  );
}

/* YouTube to AI Gradient */
.brand-gradient {
  background: linear-gradient(
    90deg,
    var(--youtube-red) 0%,
    var(--ai-purple-main) 100%
  );
}

/* Success Automation Gradient */
.auto-gradient {
  background: linear-gradient(
    135deg,
    var(--auto-green-main) 0%,
    var(--auto-green-light) 100%
  );
}

🎨 DARK MODE VARIANT
css.dark-theme {
  /* Backgrounds */
  --bg-primary: #0F0F0F;        /* YouTube dark */
  --bg-secondary: #1A1A1A;      /* Elevated */
  --bg-tertiary: #242424;       /* Cards */
  
  /* Keep brand colors vibrant */
  --youtube-red: #FF0000;
  --ai-purple-main: #8B5CF6;
  --auto-green-main: #10B981;
  
  /* Adjust text */
  --text-primary: #FAFAFA;
  --text-secondary: #A0A0A0;
}

💡 MARKA KİMLİĞİ
Logo Konsept

YouTube play button + AI sinir ağı füzyonu
Kırmızıdan mora gradient
Minimal, modern çizgiler

Typography
css--font-primary: 'Inter', -apple-system, sans-serif;
--font-mono: 'JetBrains Mono', monospace;

/* Weights */
--font-regular: 400;
--font-medium: 500;
--font-semibold: 600;
--font-bold: 700;
Spacing System
css--space-xs: 4px;
--space-sm: 8px;
--space-md: 16px;
--space-lg: 24px;
--space-xl: 32px;
--space-2xl: 48px;

🚀 ASSOS BRAND GUIDELINES
Tagline Options

"Create. Automate. Dominate."
"Your AI Video Empire"
"Content at the Speed of Thought"

Value Props in UI

YouTube Native → Kırmızı vurgu
AI Powered → Mor gradient
Fully Automated → Yeşil success indicators
Clean Interface → Off-white minimal design

 İSİM ÖNERİLERİ
Ana Öneri: ASSOS

Aristoteles'in öğrettiği antik kent
"Automated Studio System for Online Success" kısaltması da olabilir
Kolay telaffuz, global appeal
Domain: assos.ai, assos.io, assos.studio


