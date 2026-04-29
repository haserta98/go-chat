interface AvatarProps {
  name: string;
  size?: 'sm' | 'md' | 'lg';
  variant?: 'primary' | 'group';
  online?: boolean;
  children?: React.ReactNode;
}

const sizes = {
  sm: 'w-8 h-8 text-xs',
  md: 'w-9 h-9 text-sm',
  lg: 'w-10 h-10 text-base',
};

export function Avatar({ name, size = 'md', variant = 'primary', online, children }: AvatarProps) {
  const bg = variant === 'group' ? 'bg-purple-500' : 'bg-primary';

  return (
    <div className="relative shrink-0">
      <div className={`${sizes[size]} ${bg} rounded-full flex items-center justify-center font-semibold text-white`}>
        {children ?? name.charAt(0).toUpperCase()}
      </div>
      {online !== undefined && (
        <span
          className={`absolute bottom-0 right-0 w-2.5 h-2.5 rounded-full border-2 border-surface
            ${online ? 'bg-success animate-[pulse-glow_2s_ease-in-out_infinite]' : 'bg-slate-500'}`}
        />
      )}
    </div>
  );
}
