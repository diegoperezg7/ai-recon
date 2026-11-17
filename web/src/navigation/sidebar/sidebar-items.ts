import {
  Activity,
  Ban,
  BellRing,
  Bot,
  Brain,
  Bug,
  ClipboardList,
  Cpu,
  FolderOpen,
  FolderSync,
  LayoutDashboard,
  type LucideIcon,
  MessageSquare,
  Network,
  Plug,
  Radio,
  ScrollText,
  Settings2,
  Shield,
  ShieldAlert,
  Sparkles,
  Target,
  Terminal,
  Wrench,
} from "lucide-react";

export type NavBadge = "new" | "soon";

export interface NavSubItem {
  id: string;
  title: string;
  url: string;
  icon?: LucideIcon;
  badge?: NavBadge;
  disabled?: boolean;
  newTab?: boolean;
}

interface NavItemBase {
  id: string;
  title: string;
  icon?: LucideIcon;
  badge?: NavBadge;
  disabled?: boolean;
  newTab?: boolean;
}

export interface NavMainLinkItem extends NavItemBase {
  url: string;
  subItems?: never;
}

export interface NavMainParentItem extends NavItemBase {
  subItems: NavSubItem[];
}

export type NavMainItem = NavMainLinkItem | NavMainParentItem;

export interface NavGroup {
  id: number;
  label?: string;
  items: NavMainItem[];
}

export const sidebarItems: NavGroup[] = [
  {
    id: 1,
    label: "Operación",
    items: [
      { id: "dashboard", title: "Panel", url: "/dashboard", icon: LayoutDashboard },
      { id: "tasks", title: "Tareas", url: "/function/tasks", icon: Target },
      { id: "findings", title: "Hallazgos", url: "/function/findings", icon: Bug },
      { id: "chat", title: "Chat", url: "/chat", icon: MessageSquare },
    ],
  },
  {
    id: 2,
    label: "Evidencia",
    items: [
      { id: "traffic", title: "Tráfico", url: "/function/traffic", icon: Activity },
      { id: "commands", title: "Ejecución de herramientas", url: "/function/commands", icon: Terminal },
      { id: "llm-records", title: "Registros LLM", url: "/function/llm-records", icon: Radio },
      { id: "assets", title: "Activos", url: "/function/assets", icon: Network },
      { id: "sync", title: "Sincronización de activos", url: "/function/sync", icon: FolderSync },
      { id: "workspace", title: "Espacio de trabajo", url: "/function/workspace", icon: FolderOpen },
    ],
  },
  {
    id: 3,
    label: "Sistema",
    items: [
      {
        id: "ai-engine",
        title: "Motor IA",
        icon: Cpu,
        subItems: [
          { id: "llm", title: "LLM", url: "/system/llm", icon: Brain },
          { id: "agents", title: "Agentes", url: "/system/agents", icon: Bot },
          { id: "skills", title: "Skills", url: "/system/skills", icon: Sparkles },
          { id: "mcp", title: "MCP", url: "/system/mcp", icon: Plug },
          { id: "tools", title: "Herramientas", url: "/system/tools", icon: Wrench },
        ],
      },
      {
        id: "interception",
        title: "Interceptación",
        icon: Shield,
        subItems: [
          { id: "intercept", title: "Reglas de intercepción", url: "/system/intercept", icon: ShieldAlert },
          { id: "asset-intercept", title: "Bloqueo de activos", url: "/system/intercept/assets", icon: Ban },
          { id: "approvals", title: "Aprobaciones", url: "/system/intercept/approvals", icon: ClipboardList },
        ],
      },
      { id: "notify", title: "Notificaciones", url: "/system/notify", icon: BellRing },
      { id: "logs", title: "Logs", url: "/system/logs", icon: ScrollText },
      { id: "settings", title: "Configuración del sistema", url: "/system/settings", icon: Settings2 },
    ],
  },
];
