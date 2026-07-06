export type LogLevel = "debug" | "info" | "warning" | "error" | string;

export interface LiveLogEvent {
  level: LogLevel;
  message: string;
  timestamp: string;
}
