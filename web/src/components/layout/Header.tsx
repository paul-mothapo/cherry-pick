import React from "react";
import { Database, Wifi, WifiOff, AlertCircle, Slash } from "lucide-react";
import { useAppStore } from "@/stores/useAppStore";

export const Header: React.FC = () => {
  const { currentConnection } = useAppStore();

  return (
    <header className="bg-seasalt border-b border-gray-200 px-6 py-2 w-full">
      <div className="flex items-center justify-between w-full">
        <div className="flex items-center">
          <img
            src="/logos/cherry.png"
            alt="Cherry Pick"
            className="h-8 object-contain"
          />
          <span className="text-xl text-platinum mx-4">
            <Slash style={{ transform: "rotate(-30deg)" }} />
          </span>

          {currentConnection ? (
            <div className="flex items-center bg-platinum px-3 py-2 rounded-lg transition-all duration-300 hover:shadow-md">
              <div className="flex items-center">
                <div className="relative mr-2">
                  <div
                    className={`h-2 w-2 rounded-full transition-all duration-300 ${
                      currentConnection.status === "connected"
                        ? "bg-green-400 animate-pulse"
                        : currentConnection.status === "error"
                        ? "bg-turkey-red animate-pulse"
                        : "bg-yellow-400 animate-pulse"
                    }`}
                  ></div>
                  {currentConnection.status === "connected" && (
                    <div className="absolute inset-0 h-2 w-2 rounded-full bg-green-400 animate-ping opacity-75"></div>
                  )}
                </div>
                <div className="flex gap-1">
                  <span className="text-sm font-medium text-eerie-black transition-colors duration-200">
                    {currentConnection.name}
                  </span>
                  <span
                    className={`text-xs transition-colors duration-200 ${
                      currentConnection.status === "connected"
                        ? "text-green-600"
                        : currentConnection.status === "error"
                        ? "text-red-600"
                        : "text-yellow-600"
                    }`}
                  >
                    {currentConnection.status}
                  </span>
                </div>
              </div>
              <div className="ml-2 transition-all duration-300">
                {currentConnection.status === "connected" ? (
                  <Wifi className="h-4 w-4 text-green-500 animate-pulse" />
                ) : currentConnection.status === "error" ? (
                  <AlertCircle className="h-4 w-4 text-red-500 animate-bounce" />
                ) : (
                  <WifiOff className="h-4 w-4 text-yellow-500 animate-pulse" />
                )}
              </div>
            </div>
          ) : (
            <div className="flex items-center bg-platinum px-3 py-2 rounded-lg transition-all duration-300 hover:shadow-md opacity-75">
              <Database className="h-4 w-4 text-onyx mr-2 animate-pulse" />
              <span className="text-sm text-onyx">No connection selected</span>
              <WifiOff className="h-4 w-4 text-onyx ml-2 animate-pulse" />
            </div>
          )}
        </div>
      </div>
    </header>
  );
};
