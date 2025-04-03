import type { Route } from "./+types/home";
import { Card, Container, Row, Col, ProgressBar, OverlayTrigger, Tooltip } from "react-bootstrap";
import { useState, useEffect } from "react";
import { Speedometer2, Memory, HddFill, Ethernet, Thermometer, Cpu, Diagram2, Download, Upload, Database, CloudArrowDown, CloudArrowUp } from "react-bootstrap-icons";
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

export function meta({}: Route.MetaArgs) {
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
      {overview.nodes.map((node) => (
        <Card key={node.node_id} className="mb-2">
          <Card.Header className="py-2">
            <h5 className="mb-0 d-flex align-items-center justify-content-between">
              <div className="d-flex align-items-center">
                <ReactCountryFlag countryCode={node.metadata.location} className="me-2" svg />
                {node.metadata.label || node.Host.hostname}
              </div>
              <div className="d-flex align-items-center">
                <span className="me-2">{Math.floor(node.Host.uptime / 3600)}小时{Math.floor((node.Host.uptime % 3600) / 60)}分钟</span>
                <span className={`badge rounded-pill ${node.node_alive ? 'bg-success' : 'bg-danger'}`}>
                  {node.node_alive ? '在线' : '离线'}
                </span>
              </div>
            </h5>
          </Card.Header>
          <Card.Body className="py-2">
            <Row className="g-2">
              <Col md={6} lg={4}>
                <h6 className="d-flex align-items-center mb-2"><Cpu className="me-1" /> 系统资源</h6>
                <div className="small d-flex flex-wrap gap-2">
                  <OverlayTrigger placement="right" overlay={<Tooltip>CPU使用率</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <Speedometer2 className="me-1" />
                      <ProgressBar style={{width: '120px'}} now={node.percent.cpu * 100} label={`${(node.percent.cpu * 100).toFixed(1)}%`} />
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>内存使用率</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <Memory className="me-1" />
                      <ProgressBar style={{width: '120px'}} now={node.percent.mem} label={`${node.percent.mem.toFixed(1)}%`} />
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>磁盘使用率</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <HddFill className="me-1" />
                      <ProgressBar style={{width: '120px'}} now={node.percent.disk} label={`${node.percent.disk.toFixed(1)}%`} />
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>系统负载</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <Diagram2 className="me-1" />
                      <span>{node.load.load1.toFixed(2)} / {node.load.load5.toFixed(2)} / {node.load.load15.toFixed(2)}</span>
                    </div>
                  </OverlayTrigger>
                </div>
              </Col>

              <Col md={6} lg={4}>
                <h6 className="d-flex align-items-center mb-2"><Memory className="me-1" /> 内存状态</h6>
                <div className="small d-flex flex-wrap gap-2">
                  <OverlayTrigger placement="right" overlay={<Tooltip>总内存</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Database className="me-1" />{formatBytes(node.memory.total)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>已用内存</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Memory className="me-1" />{formatBytes(node.memory.used)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>可用内存</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Database className="me-1" />{formatBytes(node.memory.free)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>Swap总量</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Database className="me-1" />{formatBytes(node.swap.total)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>Swap使用</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Memory className="me-1" />{formatBytes(node.swap.used)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>Swap可用</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <Database className="me-1" />{formatBytes(node.swap.free)}
                    </div>
                  </OverlayTrigger>
                </div>
              </Col>

              <Col md={6} lg={4}>
                <h6 className="d-flex align-items-center mb-2"><HddFill className="me-1" /> 磁盘状态</h6>
                <div className="small d-flex flex-wrap gap-2">
                  <OverlayTrigger placement="right" overlay={<Tooltip>总容量</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Database className="me-1" />{formatBytes(node.disk.total)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>已用空间</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <HddFill className="me-1" />{formatBytes(node.disk.used)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>可用空间</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Database className="me-1" />{formatBytes(node.disk.free)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>读取速度</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Download className="me-1" />{formatBytes(node.disk.rx)}/s
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>写入速度</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <Upload className="me-1" />{formatBytes(node.disk.wx)}/s
                    </div>
                  </OverlayTrigger>
                </div>
              </Col>

              <Col md={6} lg={4}>
                <h6 className="d-flex align-items-center mb-2"><Ethernet className="me-1" /> 网络状态</h6>
                <div className="small d-flex flex-wrap gap-2">
                  <OverlayTrigger placement="right" overlay={<Tooltip>实时接收速度</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Download className="me-1" />{formatBytes(node.network.rx)}/s
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>实时发送速度</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <Upload className="me-1" />{formatBytes(node.network.tx)}/s
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>总接收流量</Tooltip>}>
                    <div className="d-flex align-items-center me-2">
                      <CloudArrowDown className="me-1" />{formatBytes(node.network.rb)}
                    </div>
                  </OverlayTrigger>
                  <OverlayTrigger placement="right" overlay={<Tooltip>总发送流量</Tooltip>}>
                    <div className="d-flex align-items-center">
                      <CloudArrowUp className="me-1" />{formatBytes(node.network.sb)}
                    </div>
                  </OverlayTrigger>
                </div>
              </Col>

              {node.temperature && Object.entries(node.temperature).length > 0 && (
                <Col md={6} lg={4}>
                  <h6 className="d-flex align-items-center mb-2"><Thermometer className="me-1" /> 温度信息</h6>
                  <div className="small d-flex flex-wrap gap-2">
                    {Object.entries(node.temperature).map(([sensor, temp]) => (
                      <OverlayTrigger key={sensor} placement="right" overlay={<Tooltip>{sensor}</Tooltip>}>
                        <div className="d-flex align-items-center me-2">
                          <Thermometer className="me-1" />{temp.toFixed(1)}°C
                        </div>
                      </OverlayTrigger>
                    ))}
                  </div>
                </Col>
              )}
            </Row>
          </Card.Body>
        </Card>
      ))}
    </Container>
  );
}
