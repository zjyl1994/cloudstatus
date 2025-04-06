import type { Route } from "./+types/home";
import { Container, Card, Row, Col, ProgressBar } from "react-bootstrap";
import { useState, useEffect } from "react";
import { Memory, HddFill, ArrowLeftRight, Diamond, Cpu, Hdd, Download, Upload, CloudArrowDown, CloudArrowUp, ThermometerHalf, Clock } from "react-bootstrap-icons";
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
    const interval = setInterval(fetchData, 30000);// 半分钟刷新一次
    return () => clearInterval(interval);
  }, []);

  if (!overview || overview.nodes.length === 0) {
    return <div>Loading...</div>;
  }

  return (
    <Container className="py-2">
      <Row xs={1} md={2} lg={4} className="g-3">
        {overview.nodes.map((node) => (
          <Col key={node.node_id}>
            <Card>
              <Card.Header>
                <div className="d-flex align-items-center justify-content-between">
                  <div className="d-flex align-items-center gap-2">
                    <ReactCountryFlag countryCode={node.metadata.location} svg />
                    <span>{node.metadata.label || node.Host.hostname}</span>
                  </div>
                  <span
                    className={`badge rounded-pill ${node.node_alive ? 'bg-success' : 'bg-danger'}`}
                    title={node.node_alive ? `${node.Host.uptime >= 86400 ?
                      `${Math.floor(node.Host.uptime / 86400)}天` :
                      node.Host.uptime >= 3600 ?
                        `${Math.floor(node.Host.uptime / 3600)}小时` :
                        `${Math.floor(node.Host.uptime / 60)}分钟`
                      }` : '离线'}
                  >
                    {node.node_alive ? '在线' : '离线'}
                  </span>
                </div>
              </Card.Header>
              <Card.Body>
                <div className="mb-3">
                  <div className="d-flex align-items-center gap-2 mb-1">
                    <Cpu /> <small>CPU使用率</small>
                  </div>
                  <ProgressBar
                    striped
                    now={node.percent.cpu}
                    variant={node.percent.cpu < 50 ? 'success' : node.percent.cpu < 80 ? 'warning' : 'danger'}
                    label={`${node.percent.cpu.toFixed(2)}%`}
                    title={`负载：${node.load.load1.toFixed(2)} | ${node.load.load5.toFixed(2)} | ${node.load.load15.toFixed(2)}`}
                  />
                </div>

                <div className="mb-3">
                  <div className="d-flex align-items-center gap-2 mb-1">
                    <Memory /> <small>内存使用率</small>
                  </div>
                  <ProgressBar
                    striped
                    now={node.percent.mem}
                    variant={node.percent.mem < 50 ? 'success' : node.percent.mem < 80 ? 'warning' : 'danger'}
                    label={`${node.percent.mem.toFixed(2)}%`}
                    title={`${formatBytes(node.memory.used)} / ${formatBytes(node.memory.total)}`}
                  />
                </div>

                <div className="mb-3">
                  <div className="d-flex align-items-center gap-2 mb-1">
                    <ArrowLeftRight /> <small>SWAP使用率</small>
                  </div>
                  <ProgressBar
                    striped
                    now={node.percent.swap}
                    variant={node.percent.swap < 50 ? 'success' : node.percent.swap < 80 ? 'warning' : 'danger'}
                    label={`${node.percent.swap.toFixed(2)}%`}
                    title={`${formatBytes(node.swap.used)} / ${formatBytes(node.swap.total)}`}
                  />
                </div>

                <div className="mb-3">
                  <div className="d-flex align-items-center gap-2 mb-1">
                    <HddFill /> <small>存储使用率</small>
                  </div>
                  <ProgressBar
                    striped
                    now={node.percent.disk}
                    variant={node.percent.disk < 50 ? 'success' : node.percent.disk < 80 ? 'warning' : 'danger'}
                    label={`${node.percent.disk.toFixed(2)}%`}
                    title={`${formatBytes(node.disk.used)} / ${formatBytes(node.disk.total)}`}
                  />
                </div>

                <div className="d-flex flex-column gap-2">
                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <Upload className="text-success" /> <span>上传速度</span>
                    </div>
                    <div>{formatBytes(node.network.tx)}/s</div>
                  </div>

                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <Download className="text-primary" /> <span>下载速度</span>
                    </div>
                    <div>{formatBytes(node.network.rx)}/s</div>
                  </div>


                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <Hdd className="text-success" /> <span>磁盘写入速度</span>
                    </div>
                    <div>{formatBytes(node.disk.wx)}/s</div>
                  </div>

                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <Hdd className="text-primary" /> <span>磁盘读取速度</span>
                    </div>
                    <div>{formatBytes(node.disk.rx)}/s</div>
                  </div>

                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <Clock className="text-info" /> <span>在线时长</span>
                    </div>
                    <div>{node.Host.uptime >= 86400 ?
                      `${Math.floor(node.Host.uptime / 86400)}天` :
                      node.Host.uptime >= 3600 ?
                        `${Math.floor(node.Host.uptime / 3600)}小时` :
                        `${Math.floor(node.Host.uptime / 60)}分钟`
                    }</div>
                  </div>


                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <CloudArrowUp className="text-success" /> <span>本月上传</span>
                    </div>
                    <div>{formatBytes(node.network.sb)}</div>
                  </div>

                  <div className="d-flex justify-content-between align-items-center">
                    <div className="d-flex align-items-center gap-1">
                      <CloudArrowDown className="text-primary" /> <span>本月下载</span>
                    </div>
                    <div>{formatBytes(node.network.rb)}</div>
                  </div>

                  {node.temperature && Object.entries(node.temperature).map(([sensor, temp], index) => (
                    <div key={sensor} className="d-flex justify-content-between align-items-center">
                      <div className="d-flex align-items-center gap-1">
                        <ThermometerHalf className="text-primary" /> <span>{sensor}温度</span>
                      </div>
                      <div>{temp}°C</div>
                    </div>
                  ))}
                </div>
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>
    </Container>
  );
}
