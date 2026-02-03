"use client";

export function HeatmapPreview() {
  const getLevel = (row: number, col: number) => {
    // Creating a more "realistic" looking pattern
    const val =
      Math.sin(row * 0.5) * Math.cos(col * 0.2) * 5 + ((row + col) % 5);
    return Math.abs(Math.floor(val)) % 5;
  };

  const colors = [
    "bg-muted/30",
    "bg-primary/20",
    "bg-primary/40",
    "bg-primary/70",
    "bg-primary",
  ];

  return (
    <div className="relative">
      <div className="relative flex flex-col gap-[4px] p-6 bg-card rounded-xl border shadow-sm items-center">
        <div className="flex flex-col gap-[4px] max-w-full pb-2">
          {[0, 1, 2, 3, 4, 5, 6].map((row) => (
            <div key={row} className="flex gap-[4px]">
              {Array.from({ length: 50 }).map((_, col) => {
                const level = getLevel(row, col);
                return (
                  <div
                    key={col}
                    className={`w-[11px] h-[11px] sm:w-[12px] sm:h-[12px] rounded-[2px] ${colors[level]} transition-colors hover:ring-2 hover:ring-primary/30 cursor-default`}
                  />
                );
              })}
            </div>
          ))}
        </div>

        <div className="flex items-center justify-between w-full mt-6 text-[11px] text-muted-foreground px-2 uppercase tracking-widest font-semibold flex-wrap gap-4">
          <div className="flex gap-4 sm:gap-6">
            <span>Jan</span>
            <span>Mar</span>
            <span>May</span>
            <span>Jul</span>
            <span>Sep</span>
            <span>Nov</span>
          </div>
          <div className="flex items-center gap-2">
            <span>Less</span>
            {colors.map((color, i) => (
              <div
                key={i}
                className={`w-[11px] h-[11px] rounded-[1px] ${color}`}
              />
            ))}
            <span>More</span>
          </div>
        </div>
      </div>
    </div>
  );
}
