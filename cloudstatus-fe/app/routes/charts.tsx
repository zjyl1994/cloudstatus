import type { Route } from "./+types/charts";
import { Container, Card, Row, Col } from "react-bootstrap";
import { useState, useEffect, useRef } from "react";
import * as echarts from "echarts";

const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
import { useParams } from "react-router-dom";

interface ChartsResponse {
    cpu: Array<{ time: string; value: number }>;
    memory: Array<{ time: string; value: number }>;
    swap: Array<{ time: string; value: number }>;
    disk_speed: Array<{ time: string; rx: number; tx: number }>;
    net_speed: Array<{ time: string; rx: number; tx: number }>;
    load: Array<{ time: string; load1: number; load5: number; load15: number }>;
    temperature: Record<string, Array<{ time: string; value: number }>>;
}

export function meta({ }: Route.MetaArgs) {
    return [
        { name: "description", content: "服务器监控图表" },
    ];
}

export default function Charts() {
    const [data, setData] = useState<ChartsResponse | null>(null);
    const [error, setError] = useState<string | null>(null);
    const { nodeId } = useParams();
    const cpuChartRef = useRef<HTMLDivElement>(null);
    const memoryChartRef = useRef<HTMLDivElement>(null);
    const swapChartRef = useRef<HTMLDivElement>(null);
    const diskChartRef = useRef<HTMLDivElement>(null);
    const netChartRef = useRef<HTMLDivElement>(null);
    const loadChartRef = useRef<HTMLDivElement>(null);
    const tempChartRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(`/api/charts?id=${nodeId}`);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const newData = await response.json();
                setData(newData);
                setError(null);
            } catch (error) {
                console.error('Error fetching charts data:', error);
                setError('获取图表数据失败');
            }
        };

        if (nodeId) {
            fetchData();
            const interval = setInterval(fetchData, 30000); // 半分钟更新一次
            return () => clearInterval(interval);
        }
    }, [nodeId]);

    useEffect(() => {
        if (!data) return;

        // CPU使用率图表
        if (cpuChartRef.current) {
            const chart = echarts.init(cpuChartRef.current);
            chart.setOption({
                title: { text: 'CPU使用率' },
                tooltip: { trigger: 'axis' },
                grid: { left: '10%' },
                xAxis: { type: 'category', data: data.cpu.map(item => item.time) },
                yAxis: { 
                    type: 'value', 
                    min: 0, 
                    max: 100, 
                    name: '%',
                    axisLabel: {
                        width: 50,
                        overflow: 'break'
                    }
                },
                series: [{
                    name: 'CPU',
                    type: 'line',
                    data: data.cpu.map(item => item.value),
                    areaStyle: {}
                }]
            });
        }

        // 内存使用率图表
        if (memoryChartRef.current) {
            const chart = echarts.init(memoryChartRef.current);
            chart.setOption({
                title: { text: '内存使用率' },
                tooltip: { trigger: 'axis' },
                grid: { left: '10%' },
                xAxis: { type: 'category', data: data.memory.map(item => item.time) },
                yAxis: { 
                    type: 'value', 
                    min: 0, 
                    max: 100, 
                    name: '%',
                    axisLabel: {
                        width: 50,
                        overflow: 'break'
                    }
                },
                series: [{
                    name: '内存',
                    type: 'line',
                    data: data.memory.map(item => item.value),
                    areaStyle: {}
                }]
            });
        }


        // 系统负载图表
        if (loadChartRef.current) {
            const chart = echarts.init(loadChartRef.current);
            chart.setOption({
                title: { text: '系统负载' },
                tooltip: { trigger: 'axis' },
                grid: { left: '10%' },
                xAxis: { type: 'category', data: data.load.map(item => item.time) },
                yAxis: { 
                    type: 'value',
                    axisLabel: {
                        width: 50,
                        overflow: 'break'
                    }
                },
                series: [
                    {
                        name: '1分钟',
                        type: 'line',
                        data: data.load.map(item => item.load1)
                    },
                    {
                        name: '5分钟',
                        type: 'line',
                        data: data.load.map(item => item.load5)
                    },
                    {
                        name: '15分钟',
                        type: 'line',
                        data: data.load.map(item => item.load15)
                    }
                ]
            });
        }

        // 交换分区使用率图表
        if (swapChartRef.current) {
            const chart = echarts.init(swapChartRef.current);
            chart.setOption({
                title: { text: '交换分区使用率' },
                tooltip: { trigger: 'axis' },
                grid: { left: '10%' },
                xAxis: { type: 'category', data: data.swap.map(item => item.time) },
                yAxis: { 
                    type: 'value', 
                    min: 0, 
                    max: 100, 
                    name: '%',
                    axisLabel: {
                        width: 50,
                        overflow: 'break'
                    }
                },
                series: [{
                    name: '交换分区',
                    type: 'line',
                    data: data.swap.map(item => item.value),
                    areaStyle: {}
                }]
            });
        }

        // 磁盘IO图表
        if (diskChartRef.current) {
            const chart = echarts.init(diskChartRef.current);
            chart.setOption({
                title: { text: '磁盘IO' },
                tooltip: { trigger: 'axis' },
                grid: { left: '15%' },
                xAxis: { type: 'category', data: data.disk_speed.map(item => item.time) },
                yAxis: { 
                    type: 'value',
                    axisLabel: {
                        width: 80,
                        overflow: 'break',
                        formatter: (value: number) => formatBytes(value) + '/s'
                    }
                },
                series: [
                    {
                        name: '读取',
                        type: 'line',
                        data: data.disk_speed.map(item => item.rx),
                        areaStyle: {},
                        tooltip: {
                            valueFormatter: (value: number) => formatBytes(value) + '/s'
                        }
                    },
                    {
                        name: '写入',
                        type: 'line',
                        data: data.disk_speed.map(item => item.tx),
                        areaStyle: {},
                        tooltip: {
                            valueFormatter: (value: number) => formatBytes(value) + '/s'
                        }
                    }
                ]
            });
        }

        // 网络IO图表
        if (netChartRef.current) {
            const chart = echarts.init(netChartRef.current);
            chart.setOption({
                title: { text: '网络IO' },
                tooltip: { trigger: 'axis' },
                grid: { left: '15%' },
                xAxis: { type: 'category', data: data.net_speed.map(item => item.time) },
                yAxis: { 
                    type: 'value',
                    axisLabel: {
                        width: 80,
                        overflow: 'break',
                        formatter: (value: number) => formatBytes(value) + '/s'
                    }
                },
                series: [
                    {
                        name: '接收',
                        type: 'line',
                        data: data.net_speed.map(item => item.rx),
                        areaStyle: {},
                        tooltip: {
                            valueFormatter: (value: number) => formatBytes(value) + '/s'
                        }
                    },
                    {
                        name: '发送',
                        type: 'line',
                        data: data.net_speed.map(item => item.tx),
                        areaStyle: {},
                        tooltip: {
                            valueFormatter: (value: number) => formatBytes(value) + '/s'
                        }
                    }
                ]
            });
        }

        // 温度图表
        if (tempChartRef.current && Object.keys(data.temperature).length > 0) {
            const chart = echarts.init(tempChartRef.current);
            const series = Object.entries(data.temperature).map(([name, values]) => ({
                name,
                type: 'line',
                data: values.map(item => item.value)
            }));

            chart.setOption({
                title: { text: '温度监控' },
                tooltip: { trigger: 'axis' },
                grid: { left: '10%' },
                legend: { data: Object.keys(data.temperature) },
                xAxis: {
                    type: 'category',
                    data: Object.values(data.temperature)[0].map(item => item.time)
                },
                yAxis: { 
                    type: 'value', 
                    name: '°C',
                    axisLabel: {
                        width: 50,
                        overflow: 'break'
                    }
                },
                series
            });
        }

        // 窗口大小改变时重绘图表
        const handleResize = () => {
            const charts = document.querySelectorAll('.chart-container');
            charts.forEach(container => {
                const chart = echarts.getInstanceByDom(container as HTMLElement);
                chart?.resize();
            });
        };

        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, [data]);

    if (error) {
        return <Container fluid className="py-3"><div className="alert alert-danger">{error}</div></Container>;
    }

    if (!data) {
        return <Container fluid className="py-3"><div>Loading...</div></Container>;
    }

    return (
        <Container fluid className="py-3">
            <Row>
                <Col md={6} className="mb-3">
                    <Card>
                        <Card.Body>
                            <div ref={cpuChartRef} className="chart-container" style={{ height: '300px' }} />
                        </Card.Body>
                    </Card>
                </Col>
                <Col md={6} className="mb-3">
                    <Card>
                        <Card.Body>
                            <div ref={memoryChartRef} className="chart-container" style={{ height: '300px' }} />
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
            <Row>
                <Col md={6} className="mb-3">
                    <Card>
                        <Card.Body>
                            <div ref={loadChartRef} className="chart-container" style={{ height: '300px' }} />
                        </Card.Body>
                    </Card>
                </Col>
                <Col md={6} className="mb-3">
                    <Card>
                        <Card.Body>
                            <div ref={swapChartRef} className="chart-container" style={{ height: '300px' }} />
                        </Card.Body>
                    </Card>
                </Col>

            </Row>
            <Row>
                <Col md={6} className="mb-3">
                    <Card>
                        <Card.Body>
                            <div ref={netChartRef} className="chart-container" style={{ height: '300px' }} />
                        </Card.Body>
                    </Card>
                </Col>
                <Col md={6} className="mb-3">
                    <Card>
                        <Card.Body>
                            <div ref={diskChartRef} className="chart-container" style={{ height: '300px' }} />
                        </Card.Body>
                    </Card>
                </Col>


            </Row>
            {Object.keys(data.temperature).length > 0 && (
                <Row>
                    <Col md={12} className="mb-3">
                        <Card>
                            <Card.Body>
                                <div ref={tempChartRef} className="chart-container" style={{ height: '300px' }} />
                            </Card.Body>
                        </Card>
                    </Col>
                </Row>
            )}
        </Container>
    );
}