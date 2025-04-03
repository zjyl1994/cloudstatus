import type { Route } from "./+types/home";
import { Container, Table, ProgressBar } from "react-bootstrap";
import { useState, useEffect } from "react";
import { Memory, HddFill, Ethernet, Diamond, Cpu, Hdd, Download, Upload, CloudArrowDown, CloudArrowUp } from "react-bootstrap-icons";
import ReactCountryFlag from "react-country-flag";

interface Overview {
  update_at: number;
  nodes: Array<{
    percent: {
      cpu: number;
      mem: number;
      swap: number;
      disk: number;
    };
    load: {
      load1: number;
      load5: number;
      load15: number;
    };
    memory: {
      total: number;
      used: number;
      free: number;
    };
    swap: {
      total: number;
      used: number;
      free: number;
    };
    disk: {
      total: number;
      used: number;
      free: number;
      rx: number;
      wx: number;
    };
    network: {
      rx: number;
      tx: number;
      sb: number;
      rb: number;
    };
    Host: {
      uptime: number;
      hostname: string;
      platform: string;
      version: string;
      arch: string;
    };
    node_id: string;
    interval: number;
    report: number;
    temperature: Record<string, number> | null;
    metadata: {
      id: string;
      label: string;
      location: string;
      reset_day: number;
    };
    node_alive: boolean;
  }>;
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
}

export function meta({ }: Route.MetaArgs) {
  return [
    { name: "description", content: "服务器状态监控" },
  ];
}

export default function Home() {
  const [overview, setOverview] = useState<Overview | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('/api/overview');
        const data = await response.json();
        setOverview(data);
      } catch (error) {
        console.error('Error fetching overview:', error);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 2000);
    return () => clearInterval(interval);
  }, []);

  if (!overview || overview.nodes.length === 0) {
    return <div>Loading...</div>;
  }

  return (
    <Container className="py-2">
      <Table hover responsive>
        <thead>
          <tr>
            <th>
              <div className="d-flex align-items-center gap-2">
                <Diamond /> 节点
              </div>
            </th>

            <th>
              <div className="d-flex align-items-center gap-2">
                <Cpu /> CPU 负载
              </div>
            </th>
            <th>
              <div className="d-flex align-items-center gap-2">
                <Memory /> 内存
              </div>
            </th>
            <th>
              <div className="d-flex align-items-center gap-2">
                <Hdd /> 交换
              </div>
            </th>
            <th>
              <div className="d-flex align-items-center gap-2">
                <HddFill /> 磁盘
              </div>
            </th>
            <th>
              <div className="d-flex align-items-center gap-2">
                <Ethernet /> 网络
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          {overview.nodes.map((node) => (
            <tr key={node.node_id}>
              <td>
                <div className="d-flex align-items-center gap-2">
                  <ReactCountryFlag countryCode={node.metadata.location} className="me-2" svg />
                  <span className="me-2">{node.metadata.label || node.Host.hostname}</span>
                  <span className={`badge rounded-pill ${node.node_alive ? 'bg-success' : 'bg-danger'}`} title={`${Math.floor(node.Host.uptime / 3600)}小时${Math.floor((node.Host.uptime % 3600) / 60)}分钟`}>
                    {node.node_alive ? '在线' : '离线'}
                  </span>
                </div>
              </td>
              <td>
                <div className="d-flex flex-column gap-2">
                  <ProgressBar striped now={node.percent.cpu} variant={node.percent.cpu < 50 ? 'success' : node.percent.cpu < 80 ? 'warning' : 'danger'} className="mb-0" />
                  <span className="text-nowrap">{node.load.load1.toFixed(2)} / {node.load.load5.toFixed(2)} / {node.load.load15.toFixed(2)}</span>
                </div>
              </td>
              <td>
                <div className="d-flex flex-column gap-2">
                  <ProgressBar striped now={node.percent.mem} variant={node.percent.mem < 50 ? 'success' : node.percent.mem < 80 ? 'warning' : 'danger'} className="mb-0" />
                  <span className="text-nowrap">{formatBytes(node.memory.used)}/{formatBytes(node.memory.total)}</span>
                </div>
              </td>
              <td>
                <div className="d-flex flex-column gap-2">
                  <ProgressBar striped now={node.percent.swap} variant={node.percent.swap < 50 ? 'success' : node.percent.swap < 80 ? 'warning' : 'danger'} className="mb-0" />
                  <span className="text-nowrap">{formatBytes(node.swap.used)}/{formatBytes(node.swap.total)}</span>
                </div>
              </td>
              <td>
                <div className="d-flex flex-column gap-2">
                  <ProgressBar striped now={node.percent.disk} variant={node.percent.disk < 50 ? 'success' : node.percent.disk < 80 ? 'warning' : 'danger'} className="mb-0" />
                  <span>{formatBytes(node.disk.rx)}/s {formatBytes(node.disk.wx)}/s</span>
                </div>
              </td>
              <td>
                <div className="d-flex flex-column gap-2">
                  <span>{formatBytes(node.network.rb)} {formatBytes(node.network.sb)}</span>
                  <span>{formatBytes(node.network.rx)}/s {formatBytes(node.network.tx)}/s</span>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </Container>
  );
}
