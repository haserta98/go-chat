import { Search } from 'lucide-react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  search?: boolean;
}

export function Input({ search, className = '', ...props }: InputProps) {
  if (search) {
    return (
      <div className="relative">
        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
        <input
          className={`w-full bg-black/20 border border-border text-white py-2.5 pl-9 pr-4 rounded-lg
            font-[inherit] outline-none transition-all focus:border-primary focus:ring-2 focus:ring-primary/20 ${className}`}
          {...props}
        />
      </div>
    );
  }

  return (
    <input
      className={`w-full bg-black/20 border border-border text-white py-3 px-4 rounded-lg
        font-[inherit] outline-none transition-all focus:border-primary focus:ring-2 focus:ring-primary/20 ${className}`}
      {...props}
    />
  );
}
