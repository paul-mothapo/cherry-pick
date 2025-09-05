import React from "react";
import { NavLink } from "react-router-dom";
import {
  Database,
  // Shield,
  BarChart3,
  // Search,
  // AlertTriangle,
  // GitBranch,
  // Settings,
  LayoutDashboard,
  PanelLeftClose,
  PanelLeftOpen,
} from "lucide-react";
import { clsx } from "clsx";
import { useAppStore } from "@/stores/useAppStore";
import { Button } from "@/components/ui/Button";

const navigation = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Connections", href: "/connections", icon: Database },
  { name: "Analysis", href: "/analysis", icon: BarChart3 },
  // { name: "Security", href: "/security", icon: Shield },
  // { name: "Optimization", href: "/optimization", icon: Search },
  // { name: "Monitoring", href: "/monitoring", icon: AlertTriangle },
  // { name: "Lineage", href: "/lineage", icon: GitBranch },
  // { name: "Settings", href: "/settings", icon: Settings },
];

export const Sidebar: React.FC = () => {
  const { sidebarOpen, sidebarCollapsed, setSidebarCollapsed } = useAppStore();

  if (!sidebarOpen) return null;

  return (
    <div
      className={clsx(
        "flex flex-col bg-seasalt border-r border-gray-200 h-full transition-all duration-300 ease-in-out",
        sidebarCollapsed ? "w-16" : "w-60"
      )}
    >
      <nav className="flex-1 px-4 py-6 space-y-1 mt-4">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            title={sidebarCollapsed ? item.name : undefined}
            className={({ isActive }) =>
              clsx(
                "group flex items-center text-sm font-medium rounded-md transition-colors duration-200",
                sidebarCollapsed ? "px-3 py-3 justify-center" : "px-3 py-2",
                isActive
                  ? "bg-platinum text-eerie-black"
                  : "text-onyx hover:bg-platinum hover:text-onyx"
              )
            }
          >
            <item.icon
              className={clsx(
                "h-5 w-5 flex-shrink-0",
                sidebarCollapsed ? "mr-0" : "mr-3"
              )}
            />
            {!sidebarCollapsed && (
              <span className="transition-opacity duration-200">
                {item.name}
              </span>
            )}
          </NavLink>
        ))}
      </nav>

      <div
        className={clsx(
          "px-4 py-4 border-t border-gray-200 flex items-center",
          sidebarCollapsed ? "justify-center" : "justify-between"
        )}
      >
        {!sidebarCollapsed && (
          <div className="text-xs text-onyx">Cherry Pick v1.0</div>
        )}
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
          title={sidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
          className={clsx(sidebarCollapsed ? "mx-auto" : "")}
        >
          {sidebarCollapsed ? (
            <PanelLeftOpen className="h-4 w-4" />
          ) : (
            <PanelLeftClose className="h-4 w-4" />
          )}
        </Button>
      </div>
    </div>
  );
};
