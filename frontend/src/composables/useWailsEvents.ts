import { Events } from "@wailsio/runtime";
import { useLogsStore } from "@/stores/logs";
import type { LiveLogEvent } from "@/types/log";

let registered = false;

export function registerWailsEvents() {
  if (registered) {
    return;
  }
  registered = true;

  const logs = useLogsStore();
  Events.On("log", (event: { data: LiveLogEvent }) => {
    logs.addLog(event.data);
  });
}
