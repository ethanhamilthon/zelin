import { call } from './bridge';

export type MCPServer = { name: string; command: string; args: string[]; env?: Record<string, string> };
export type MCPTool = { name: string; description?: string };

export const listServers = () => call<MCPServer[]>('MCPService.List');
export const addServer = (server: MCPServer) => call<MCPServer[]>('MCPService.Add', server);
export const removeServer = (name: string) => call<MCPServer[]>('MCPService.Remove', name);
export const testServer = (name: string) => call<MCPTool[]>('MCPService.Test', name);

export function parseServer(name: string, line: string, envText: string): MCPServer {
  const parts = line.match(/"[^"]*"|'[^']*'|\S+/g)?.map((part) => part.replace(/^(["'])(.*)\1$/, '$2')) ?? [];
  const env: Record<string, string> = {};
  for (const row of envText.split('\n')) {
    const at = row.indexOf('=');
    if (at > 0) env[row.slice(0, at).trim()] = row.slice(at + 1).trim();
  }
  return { name: name.trim(), command: parts[0] ?? '', args: parts.slice(1), ...(Object.keys(env).length && { env }) };
}
