import { Loader2 } from 'lucide-react';

type Variant = 'primary' | 'ghost' | 'danger';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  loading?: boolean;
  icon?: React.ReactNode;
}

const base = 'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed';

const variants: Record<Variant, string> = {
  primary: 'bg-primary text-white hover:bg-primary-hover active:scale-[0.98] px-4 py-2.5',
  ghost: 'bg-transparent text-slate-300 hover:bg-white/5 px-2 py-2',
  danger: 'bg-red-500/20 text-red-400 border border-red-500/30 hover:bg-red-500/30 px-4 py-2.5',
};

export function Button({ variant = 'primary', loading, icon, children, className = '', ...props }: ButtonProps) {
  return (
    <button className={`${base} ${variants[variant]} ${className}`} disabled={loading || props.disabled} {...props}>
      {loading ? <Loader2 size={16} className="animate-spin" /> : icon}
      {children}
    </button>
  );
}

interface IconButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  loading?: boolean;
  color?: 'success' | 'danger' | 'dangerHidden' | 'default';
}

const iconColors = {
  success: 'text-success hover:bg-success/15',
  danger: 'text-danger hover:bg-danger/15',
  dangerHidden: 'text-danger hover:bg-danger/15 opacity-0 group-hover:opacity-100',
  default: 'text-slate-400 hover:bg-white/10',
};

export function IconButton({ loading, color = 'default', children, className = '', ...props }: IconButtonProps) {
  return (
    <button
      className={`inline-flex items-center justify-center w-7 h-7 rounded-lg bg-transparent border-none
        cursor-pointer transition-all shrink-0 disabled:opacity-50 disabled:cursor-not-allowed
        hover:scale-110 ${iconColors[color]} ${className}`}
      disabled={loading || props.disabled}
      {...props}
    >
      {loading ? <Loader2 size={14} className="animate-spin" /> : children}
    </button>
  );
}
