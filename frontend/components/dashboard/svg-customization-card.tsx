"use client";

import { Palette } from "lucide-react";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Slider } from "@/components/ui/slider";
import { Switch } from "@/components/ui/switch";
import { SVGOptions } from "@/lib/schemas";

interface SVGCustomizationCardProps {
  options: SVGOptions;
  setOptions: (options: SVGOptions) => void;
  themes: Array<{ id: string; name: string }>;
}

export function SVGCustomizationCard({
  options,
  setOptions,
  themes,
}: SVGCustomizationCardProps) {
  const updateOption = <K extends keyof SVGOptions>(
    key: K,
    value: SVGOptions[K],
  ) => {
    setOptions({ ...options, [key]: value });
  };

  return (
    <Card>
      <CardHeader className="pb-4">
        <CardTitle className="text-base flex items-center gap-2">
          <Palette className="h-4 w-4" />
          Appearance
        </CardTitle>
        <CardDescription>Choose a theme and adjust the look</CardDescription>
      </CardHeader>
      <CardContent className="space-y-5">
        {/* Theme */}
        <div className="space-y-2">
          <Label className="text-sm font-medium">Color Theme</Label>
          <Select
            value={options.theme}
            onValueChange={(value) => updateOption("theme", value)}
          >
            <SelectTrigger className="h-9">
              <SelectValue placeholder="Select theme" />
            </SelectTrigger>
            <SelectContent>
              {themes.map((theme) => (
                <SelectItem key={theme.id} value={theme.id}>
                  {theme.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Days */}
        <div className="space-y-3">
          <div className="flex justify-between">
            <Label className="text-sm font-medium">Time Range</Label>
            <span className="text-sm text-muted-foreground">
              {options.days} days
            </span>
          </div>
          <Slider
            value={[options.days || 365]}
            onValueChange={([value]) => updateOption("days", value)}
            min={30}
            max={365}
            step={30}
          />
        </div>

        {/* Cell Size */}
        <div className="space-y-3">
          <div className="flex justify-between">
            <Label className="text-sm font-medium">Cell Size</Label>
            <span className="text-sm text-muted-foreground">
              {options.cell_size}px
            </span>
          </div>
          <Slider
            value={[options.cell_size || 11]}
            onValueChange={([value]) => updateOption("cell_size", value)}
            min={5}
            max={20}
            step={1}
          />
        </div>

        {/* Border Radius */}
        <div className="space-y-3">
          <div className="flex justify-between">
            <Label className="text-sm font-medium">Corner Radius</Label>
            <span className="text-sm text-muted-foreground">
              {options.radius}px
            </span>
          </div>
          <Slider
            value={[options.radius ?? 2]}
            onValueChange={([value]) => updateOption("radius", value)}
            min={0}
            max={10}
            step={1}
          />
        </div>

        {/* Divider */}
        <div className="border-t pt-4">
          <Label className="text-sm font-medium mb-3 block">
            Display Options
          </Label>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm">Show Legend</span>
              <Switch
                checked={!options.hide_legend}
                onCheckedChange={(checked) =>
                  updateOption("hide_legend", !checked)
                }
              />
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm">Show Total Count</span>
              <Switch
                checked={!options.hide_total}
                onCheckedChange={(checked) =>
                  updateOption("hide_total", !checked)
                }
              />
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm">Show Day/Month Labels</span>
              <Switch
                checked={!options.hide_labels}
                onCheckedChange={(checked) =>
                  updateOption("hide_labels", !checked)
                }
              />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
