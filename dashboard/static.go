package dashboard

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// TODO: use https://github.com/rakyll/statik
// Right now, only fontawesome is not working.
func staticHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	fmt.Fprintln(w, htmlSourceFile)
	return nil
}

var htmlSourceFile = `
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="author" content="Gyu-Ho Lee">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/tether/1.1.1/js/tether.min.js"></script>
    <script src="https://code.jquery.com/jquery-2.2.0.min.js"></script>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.2/css/bootstrap.min.css" integrity="sha384-y3tfxAZXuh4HwSYylfB+J125MxIs6mR5FOHamPBG064zB+AFeWH94NdvaCBm8qnd" crossorigin="anonymous">
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.2/js/bootstrap.min.js" integrity="sha384-vZ2WRJMwsjRMW/8U7i6PWi6AlO1L79snBrmgiDpgIWJ82z8eA5lenwvxbMV1PAh7" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css">
    <script src="https://code.highcharts.com/highcharts.js"></script>
    <script src="https://code.highcharts.com/modules/exporting.js"></script>
    <script>
    $(document).ready(function() {
        $('#endpoint').submit(function(e) {
            e.preventDefault();
            $.ajax({
                type: "POST",
                url: "/endpoint",
                data: $(this).serialize(),
                success: function(dataObj) {
                    $.ajax({
                        type: "GET",
                        url: "/endpoint",
                        async: true,
                        dataType: "json",
                        success: function(dataObj) {
                            var refreshRate = dataObj.RefreshInMillisecond;

                            if (dataObj.Success == true) {
                                document.getElementById('endpoint_state_button').className = "btn btn-primary";
                                document.getElementById('endpoint_state').className = "fa fa-check";
                            } else {
                                document.getElementById('endpoint_state_button').className = "btn btn-danger";
                                document.getElementById('endpoint_state').className = "fa fa-close";
                            }

                            document.getElementById('endpoint_result').innerHTML = dataObj.Message;

                            if (dataObj.Success == true) {

                                setInterval(function() {
                                    $.ajax({
                                        type: "GET",
                                        url: '/stats',
                                        async: true,
                                        dataType: "json",
                                        success: function(dataObj) {
                                            document.getElementById('etcd1_Name').innerHTML = "Name: <b>" + dataObj.Etcd1Name + "</b>";
                                            document.getElementById('etcd1_Endpoint').innerHTML = "Endpoint: <b>" + dataObj.Etcd1Endpoint + "</b>";
                                            document.getElementById('etcd1_ID').innerHTML = "ID: <b>" + dataObj.Etcd1ID + "</b>";
                                            document.getElementById('etcd1_StartTime').innerHTML = "StartTime: <b>" + dataObj.Etcd1StartTime + "</b>";
                                            document.getElementById('etcd1_LeaderUptime').innerHTML = "LeaderUptime: <b>" + dataObj.Etcd1LeaderUptime + "</b>";
                                            document.getElementById('etcd1_RecvAppendRequestCnt').innerHTML = "RecvAppendRequestCnt: <b>" + dataObj.Etcd1RecvAppendRequestCnt + "</b>";
                                            document.getElementById('etcd1_RecvingBandwidthRate').innerHTML = "RecvingBandwidthRate: <b>" + dataObj.Etcd1RecvingBandwidthRate + "</b>";
                                            document.getElementById('etcd1_SendAppendRequestCnt').innerHTML = "SendAppendRequestCnt: <b>" + dataObj.Etcd1SendAppendRequestCnt + "</b>";
                                            document.getElementById('etcd1_SendingBandwidthRate').innerHTML = "SendingBandwidthRate: <b>" + dataObj.Etcd1SendingBandwidthRate + "</b>";
                                            if (dataObj.Etcd1State == "StateLeader") {
                                                document.getElementById('etcd1_state_button').className = "btn btn-primary";
                                                document.getElementById('etcd1_state').className = "fa fa-fort-awesome";
                                            } else if (dataObj.Etcd1State == "StateCandidate") {
                                                document.getElementById('etcd1_state_button').className = "btn btn-info";
                                                document.getElementById('etcd1_state').className = "fa fa-star-half-o";
                                            } else if (dataObj.Etcd1State == "StateFollower") {
                                                document.getElementById('etcd1_state_button').className = "btn btn-warning";
                                                document.getElementById('etcd1_state').className = "fa fa-child";
                                            } else if (dataObj.Etcd1State == "Unreachable") {
                                                document.getElementById('etcd1_state_button').className = "btn btn-danger";
                                                document.getElementById('etcd1_state').className = "fa fa-frown-o";
                                            } else {
                                                document.getElementById('etcd1_state_button').className = "btn btn-secondary-outline";
                                                document.getElementById('etcd1_state').className = "fa fa-pause-circle";
                                            }

                                            document.getElementById('etcd2_Name').innerHTML = "Name: <b>" + dataObj.Etcd2Name + "</b>";
                                            document.getElementById('etcd2_Endpoint').innerHTML = "Endpoint: <b>" + dataObj.Etcd2Endpoint + "</b>";
                                            document.getElementById('etcd2_ID').innerHTML = "ID: <b>" + dataObj.Etcd2ID + "</b>";
                                            document.getElementById('etcd2_StartTime').innerHTML = "StartTime: <b>" + dataObj.Etcd2StartTime + "</b>";
                                            document.getElementById('etcd2_LeaderUptime').innerHTML = "LeaderUptime: <b>" + dataObj.Etcd2LeaderUptime + "</b>";
                                            document.getElementById('etcd2_RecvAppendRequestCnt').innerHTML = "RecvAppendRequestCnt: <b>" + dataObj.Etcd2RecvAppendRequestCnt + "</b>";
                                            document.getElementById('etcd2_RecvingBandwidthRate').innerHTML = "RecvingBandwidthRate: <b>" + dataObj.Etcd2RecvingBandwidthRate + "</b>";
                                            document.getElementById('etcd2_SendAppendRequestCnt').innerHTML = "SendAppendRequestCnt: <b>" + dataObj.Etcd2SendAppendRequestCnt + "</b>";
                                            document.getElementById('etcd2_SendingBandwidthRate').innerHTML = "SendingBandwidthRate: <b>" + dataObj.Etcd2SendingBandwidthRate + "</b>";
                                            if (dataObj.Etcd2State == "StateLeader") {
                                                document.getElementById('etcd2_state_button').className = "btn btn-primary";
                                                document.getElementById('etcd2_state').className = "fa fa-fort-awesome";
                                            } else if (dataObj.Etcd2State == "StateCandidate") {
                                                document.getElementById('etcd2_state_button').className = "btn btn-info";
                                                document.getElementById('etcd2_state').className = "fa fa-star-half-o";
                                            } else if (dataObj.Etcd2State == "StateFollower") {
                                                document.getElementById('etcd2_state_button').className = "btn btn-warning";
                                                document.getElementById('etcd2_state').className = "fa fa-child";
                                            } else if (dataObj.Etcd2State == "Unreachable") {
                                                document.getElementById('etcd2_state_button').className = "btn btn-danger";
                                                document.getElementById('etcd2_state').className = "fa fa-frown-o";
                                            } else {
                                                document.getElementById('etcd2_state_button').className = "btn btn-secondary-outline";
                                                document.getElementById('etcd2_state').className = "fa fa-pause-circle";
                                            }

                                            document.getElementById('etcd3_Name').innerHTML = "Name: <b>" + dataObj.Etcd3Name + "</b>";
                                            document.getElementById('etcd3_Endpoint').innerHTML = "Endpoint: <b>" + dataObj.Etcd3Endpoint + "</b>";
                                            document.getElementById('etcd3_ID').innerHTML = "ID: <b>" + dataObj.Etcd3ID + "</b>";
                                            document.getElementById('etcd3_StartTime').innerHTML = "StartTime: <b>" + dataObj.Etcd3StartTime + "</b>";
                                            document.getElementById('etcd3_LeaderUptime').innerHTML = "LeaderUptime: <b>" + dataObj.Etcd3LeaderUptime + "</b>";
                                            document.getElementById('etcd3_RecvAppendRequestCnt').innerHTML = "RecvAppendRequestCnt: <b>" + dataObj.Etcd3RecvAppendRequestCnt + "</b>";
                                            document.getElementById('etcd3_RecvingBandwidthRate').innerHTML = "RecvingBandwidthRate: <b>" + dataObj.Etcd3RecvingBandwidthRate + "</b>";
                                            document.getElementById('etcd3_SendAppendRequestCnt').innerHTML = "SendAppendRequestCnt: <b>" + dataObj.Etcd3SendAppendRequestCnt + "</b>";
                                            document.getElementById('etcd3_SendingBandwidthRate').innerHTML = "SendingBandwidthRate: <b>" + dataObj.Etcd3SendingBandwidthRate + "</b>";
                                            if (dataObj.Etcd3State == "StateLeader") {
                                                document.getElementById('etcd3_state_button').className = "btn btn-primary";
                                                document.getElementById('etcd3_state').className = "fa fa-fort-awesome";
                                            } else if (dataObj.Etcd3State == "StateCandidate") {
                                                document.getElementById('etcd3_state_button').className = "btn btn-info";
                                                document.getElementById('etcd3_state').className = "fa fa-star-half-o";
                                            } else if (dataObj.Etcd3State == "StateFollower") {
                                                document.getElementById('etcd3_state_button').className = "btn btn-warning";
                                                document.getElementById('etcd3_state').className = "fa fa-child";
                                            } else if (dataObj.Etcd3State == "Unreachable") {
                                                document.getElementById('etcd3_state_button').className = "btn btn-danger";
                                                document.getElementById('etcd3_state').className = "fa fa-frown-o";
                                            } else {
                                                document.getElementById('etcd3_state_button').className = "btn btn-secondary-outline";
                                                document.getElementById('etcd3_state').className = "fa fa-pause-circle";
                                            }

                                            document.getElementById('etcd4_Endpoint').innerHTML = "Endpoint: <b>" + dataObj.Etcd4Endpoint + "</b>";
                                            document.getElementById('etcd4_Name').innerHTML = "Name: <b>" + dataObj.Etcd4Name + "</b>";
                                            document.getElementById('etcd4_ID').innerHTML = "ID: <b>" + dataObj.Etcd4ID + "</b>";
                                            document.getElementById('etcd4_StartTime').innerHTML = "StartTime: <b>" + dataObj.Etcd4StartTime + "</b>";
                                            document.getElementById('etcd4_LeaderUptime').innerHTML = "LeaderUptime: <b>" + dataObj.Etcd4LeaderUptime + "</b>";
                                            document.getElementById('etcd4_RecvAppendRequestCnt').innerHTML = "RecvAppendRequestCnt: <b>" + dataObj.Etcd4RecvAppendRequestCnt + "</b>";
                                            document.getElementById('etcd4_RecvingBandwidthRate').innerHTML = "RecvingBandwidthRate: <b>" + dataObj.Etcd4RecvingBandwidthRate + "</b>";
                                            document.getElementById('etcd4_SendAppendRequestCnt').innerHTML = "SendAppendRequestCnt: <b>" + dataObj.Etcd4SendAppendRequestCnt + "</b>";
                                            document.getElementById('etcd4_SendingBandwidthRate').innerHTML = "SendingBandwidthRate: <b>" + dataObj.Etcd4SendingBandwidthRate + "</b>";
                                            if (dataObj.Etcd4State == "StateLeader") {
                                                document.getElementById('etcd4_state_button').className = "btn btn-primary";
                                                document.getElementById('etcd4_state').className = "fa fa-fort-awesome";
                                            } else if (dataObj.Etcd4State == "StateCandidate") {
                                                document.getElementById('etcd4_state_button').className = "btn btn-info";
                                                document.getElementById('etcd4_state').className = "fa fa-star-half-o";
                                            } else if (dataObj.Etcd4State == "StateFollower") {
                                                document.getElementById('etcd4_state_button').className = "btn btn-warning";
                                                document.getElementById('etcd4_state').className = "fa fa-child";
                                            } else if (dataObj.Etcd4State == "Unreachable") {
                                                document.getElementById('etcd4_state_button').className = "btn btn-danger";
                                                document.getElementById('etcd4_state').className = "fa fa-frown-o";
                                            } else {
                                                document.getElementById('etcd4_state_button').className = "btn btn-secondary-outline";
                                                document.getElementById('etcd4_state').className = "fa fa-pause-circle";
                                            }

                                            document.getElementById('etcd5_Endpoint').innerHTML = "Endpoint: <b>" + dataObj.Etcd5Endpoint + "</b>";
                                            document.getElementById('etcd5_Name').innerHTML = "Name: <b>" + dataObj.Etcd5Name + "</b>";
                                            document.getElementById('etcd5_ID').innerHTML = "ID: <b>" + dataObj.Etcd5ID + "</b>";
                                            document.getElementById('etcd5_StartTime').innerHTML = "StartTime: <b>" + dataObj.Etcd5StartTime + "</b>";
                                            document.getElementById('etcd5_LeaderUptime').innerHTML = "LeaderUptime: <b>" + dataObj.Etcd5LeaderUptime + "</b>";
                                            document.getElementById('etcd5_RecvAppendRequestCnt').innerHTML = "RecvAppendRequestCnt: <b>" + dataObj.Etcd5RecvAppendRequestCnt + "</b>";
                                            document.getElementById('etcd5_RecvingBandwidthRate').innerHTML = "RecvingBandwidthRate: <b>" + dataObj.Etcd5RecvingBandwidthRate + "</b>";
                                            document.getElementById('etcd5_SendAppendRequestCnt').innerHTML = "SendAppendRequestCnt: <b>" + dataObj.Etcd5SendAppendRequestCnt + "</b>";
                                            document.getElementById('etcd5_SendingBandwidthRate').innerHTML = "SendingBandwidthRate: <b>" + dataObj.Etcd5SendingBandwidthRate + "</b>";
                                            if (dataObj.Etcd5State == "StateLeader") {
                                                document.getElementById('etcd5_state_button').className = "btn btn-primary";
                                                document.getElementById('etcd5_state').className = "fa fa-fort-awesome";
                                            } else if (dataObj.Etcd5State == "StateCandidate") {
                                                document.getElementById('etcd5_state_button').className = "btn btn-info";
                                                document.getElementById('etcd5_state').className = "fa fa-star-half-o";
                                            } else if (dataObj.Etcd5State == "StateFollower") {
                                                document.getElementById('etcd5_state_button').className = "btn btn-warning";
                                                document.getElementById('etcd5_state').className = "fa fa-child";
                                            } else if (dataObj.Etcd5State == "Unreachable") {
                                                document.getElementById('etcd5_state_button').className = "btn btn-danger";
                                                document.getElementById('etcd5_state').className = "fa fa-frown-o";
                                            } else {
                                                document.getElementById('etcd5_state_button').className = "btn btn-secondary-outline";
                                                document.getElementById('etcd5_state').className = "fa fa-pause-circle";
                                            }
                                        }
                                    });
                                }, refreshRate);

                                Highcharts.setOptions({
                                    global: {
                                        useUTC: false
                                    }
                                });

                                var fasterRefresh = refreshRate - 500;
                                var metricsQueue = new Array;
                                setInterval(function() {
                                    $.ajax({
                                        type: "GET",
                                        url: '/metrics',
                                        async: true,
                                        dataType: "json",
                                        success: function(dataObj) {
                                            metricsQueue.push(dataObj);
                                        }
                                    });
                                }, fasterRefresh);

                                $('#chart_StorageKeysTotal').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1StorageKeysTotal,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2StorageKeysTotal,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3StorageKeysTotal,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4StorageKeysTotal,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5StorageKeysTotal,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Number of keys'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Count'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_StorageBytes').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1StorageBytes,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint,
                                                            bytes: dataObj.Etcd1StorageBytesStr
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2StorageBytes,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint,
                                                            bytes: dataObj.Etcd2StorageBytesStr
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3StorageBytes,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint,
                                                            bytes: dataObj.Etcd3StorageBytesStr
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4StorageBytes,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint,
                                                            bytes: dataObj.Etcd4StorageBytesStr
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5StorageBytes,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint,
                                                            bytes: dataObj.Etcd5StorageBytesStr
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Storage size'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Bytes'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2) + " (" + this.point.bytes + ")";
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_WalFsyncSecondsSum').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1WalFsyncSecondsSum,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2WalFsyncSecondsSum,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3WalFsyncSecondsSum,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4WalFsyncSecondsSum,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5WalFsyncSecondsSum,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'WAL Fsync Seconds'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Seconds'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });


                                $('#chart_GcSecondsSum').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1GcSecondsSum,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2GcSecondsSum,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3GcSecondsSum,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4GcSecondsSum,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5GcSecondsSum,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'GC Seconds'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Seconds'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_MemstatsAllocBytes').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1MemstatsAllocBytes,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint,
                                                            bytes: dataObj.Etcd1MemstatsAllocBytesStr
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2MemstatsAllocBytes,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint,
                                                            bytes: dataObj.Etcd2MemstatsAllocBytesStr
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3MemstatsAllocBytes,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint,
                                                            bytes: dataObj.Etcd3MemstatsAllocBytesStr
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4MemstatsAllocBytes,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint,
                                                            bytes: dataObj.Etcd4MemstatsAllocBytesStr
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5MemstatsAllocBytes,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint,
                                                            bytes: dataObj.Etcd5MemstatsAllocBytesStr
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Memstats Alloc Bytes'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Bytes'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2) + " (" + this.point.bytes + ")";
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_MemstatsHeapAllocBytes').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1MemstatsHeapAllocBytes,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint,
                                                            bytes: dataObj.Etcd1MemstatsHeapAllocBytesStr
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2MemstatsHeapAllocBytes,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint,
                                                            bytes: dataObj.Etcd2MemstatsHeapAllocBytesStr
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3MemstatsHeapAllocBytes,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint,
                                                            bytes: dataObj.Etcd3MemstatsHeapAllocBytesStr
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4MemstatsHeapAllocBytes,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint,
                                                            bytes: dataObj.Etcd4MemstatsHeapAllocBytesStr
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5MemstatsHeapAllocBytes,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint,
                                                            bytes: dataObj.Etcd5MemstatsHeapAllocBytesStr
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Memstats Heap Alloc Bytes'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Bytes'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2) + " (" + this.point.bytes + ")";
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0,
                                                    bytes: "0 bytes"
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_MemstatsMallocsTotal').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1MemstatsMallocsTotal,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2MemstatsMallocsTotal,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3MemstatsMallocsTotal,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4MemstatsMallocsTotal,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5MemstatsMallocsTotal,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Memstats Mallocs'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Count'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_ProcessCPUSeconds').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1ProcessCPUSeconds,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2ProcessCPUSeconds,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3ProcessCPUSeconds,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4ProcessCPUSeconds,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5ProcessCPUSeconds,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Process CPU seconds'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Seconds'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_Goroutines').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1Goroutines,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2Goroutines,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3Goroutines,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4Goroutines,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5Goroutines,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Goroutines'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Count'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_StorageWatcherTotal').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1StorageWatcherTotal,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2StorageWatcherTotal,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3StorageWatcherTotal,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4StorageWatcherTotal,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5StorageWatcherTotal,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Storage Watchers'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Count'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_StorageWatchStreamTotal').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1StorageWatchStreamTotal,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2StorageWatchStreamTotal,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3StorageWatchStreamTotal,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4StorageWatchStreamTotal,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5StorageWatchStreamTotal,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Storage Watch Stream'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Count'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                                $('#chart_StorageSlowWatcherTotal').highcharts({
                                    chart: {
                                        type: 'spline',
                                        animation: Highcharts.svg, // don't animate in old IE
                                        marginRight: 10,
                                        events: {
                                            load: function() {
                                                var series1 = this.series[0];
                                                var series2 = this.series[1];
                                                var series3 = this.series[2];
                                                var series4 = this.series[3];
                                                var series5 = this.series[4];
                                                setInterval(function() {
                                                    var copiedQueue = metricsQueue.slice()
                                                    if (copiedQueue.length > 0) {
                                                        var dataObj = copiedQueue.shift();
                                                        var x = (new Date()).getTime();
                                                        series1.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd1StorageSlowWatcherTotal,
                                                            name: dataObj.Etcd1Name,
                                                            endpoint: dataObj.Etcd1Endpoint
                                                        }, true, true);
                                                        series2.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd2StorageSlowWatcherTotal,
                                                            name: dataObj.Etcd2Name,
                                                            endpoint: dataObj.Etcd2Endpoint
                                                        }, true, true);
                                                        series3.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd3StorageSlowWatcherTotal,
                                                            name: dataObj.Etcd3Name,
                                                            endpoint: dataObj.Etcd3Endpoint
                                                        }, true, true);
                                                        series4.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd4StorageSlowWatcherTotal,
                                                            name: dataObj.Etcd4Name,
                                                            endpoint: dataObj.Etcd4Endpoint
                                                        }, true, true);
                                                        series5.addPoint({
                                                            x: x,
                                                            y: dataObj.Etcd5StorageSlowWatcherTotal,
                                                            name: dataObj.Etcd5Name,
                                                            endpoint: dataObj.Etcd5Endpoint
                                                        }, true, true);
                                                    }
                                                }, refreshRate);
                                            }
                                        }
                                    },
                                    title: {
                                        text: 'Storage Slow Watchers'
                                    },
                                    xAxis: {
                                        type: 'datetime',
                                        tickPixelInterval: 150
                                    },
                                    yAxis: {
                                        title: {
                                            text: 'Count'
                                        },
                                        plotLines: [{
                                            value: 0,
                                            width: 1,
                                            color: '#808080'
                                        }]
                                    },
                                    tooltip: {
                                        formatter: function() {
                                            return '<b>' + this.point.name + "(" + this.point.endpoint + ")" + '</b><br/>' +
                                                Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                                                Highcharts.numberFormat(this.y, 2);
                                        }
                                    },
                                    legend: {
                                        enabled: false
                                    },
                                    exporting: {
                                        enabled: false
                                    },
                                    series: [{
                                        name: 'etcd1',
                                        color: "#FFCC00", // orange
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd2',
                                        color: "#04B4E3", // blue
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd3',
                                        color: "#00695C", // teal
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd4',
                                        color: "#CCCCCC", // gray
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }, {
                                        name: 'etcd5',
                                        color: "#C62828", // red
                                        data: (function() {
                                            var data = [],
                                                time = (new Date()).getTime(),
                                                i;
                                            for (i = -19; i <= 0; i += 1) {
                                                data.push({
                                                    x: time + i * refreshRate,
                                                    y: 0
                                                });
                                            }
                                            return data;
                                        }())
                                    }]
                                });

                            }
                        }
                    });
                }
            });
        });
    });
    </script>
    <style>
    /*
 * Base structure
 */
    /* Move down content because we have a fixed navbar that is 50px tall */
    
    body {
        padding-top: 50px;
    }
    /*
 * Global add-ons
 */
    
    .sub-header {
        padding-bottom: 10px;
        border-bottom: 1px solid #eee;
    }
    /*
 * Top navigation
 * Hide default border to remove 1px line.
 */
    
    .navbar-fixed-top {
        border: 0;
    }
    /*
 * Sidebar
 */
    /* Hide for mobile, show later */
    
    .sidebar {
        display: none;
    }
    
    @media (min-width: 350px) {
        .sidebar {
            width: 120px;
            position: fixed;
            top: 51px;
            bottom: 0;
            left: 0;
            z-index: 1000;
            display: block;
            padding: 20px;
            overflow-x: hidden;
            overflow-y: auto;
            /* Scrollable contents if viewport is shorter than content. */
            background-color: #f5f5f5;
            border-right: 1px solid #eee;
        }
    }
    /* Sidebar navigation */
    
    .nav-sidebar {
        margin-right: -21px;
        /* 20px padding + 1px border */
        margin-bottom: 20px;
        margin-left: -20px;
    }
    
    .nav-sidebar > li > a {
        padding-right: 20px;
        padding-left: 20px;
    }
    
    .nav-sidebar > .active > a,
    .nav-sidebar > .active > a:hover,
    .nav-sidebar > .active > a:focus {
        color: #fff;
        background-color: #428bca;
    }
    /*
 * Main content
 */
    
    .main {
        padding: 20px;
    }
    
    @media (min-width: 768px) {
        .main {
            padding-right: 20px;
            padding-left: 5px;
        }
    }
    
    .main .page-header {
        margin-top: 0;
    }
    /*
 * Placeholder dashboard ideas
 */
    
    .placeholders {
        margin-bottom: 30px;
        text-align: center;
    }
    
    .placeholders h4 {
        margin-bottom: 0;
    }
    
    .placeholder {
        margin-bottom: 20px;
    }
    
    .placeholder img {
        display: inline-block;
        border-radius: 50%;
    }
    
    .wrapper {
        text-align: center;
    }
    
    .ui-menu-item {
        width: 100%;
        height: 100%;
        font-family: "Lucida Console", Courier, monospace;
        font-size: 12px;
    }
    </style>
    <title>runetcd dashboard</title>
</head>

<body>
    <nav class="navbar navbar-dark navbar-fixed-top bg-inverse">
        <button type="button" class="navbar-toggler hidden-sm-up" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
        </button>
        <a class="navbar-brand" href="https://github.com/gophergala2016/runetcd">ETCD Dashboard</a>
        <div id="navbar">
            <nav class="nav navbar-nav pull-xs-left">
                <a class="nav-item nav-link" href="https://github.com/coreos/etcd"><i class="fa fa-github fa-2x"></i></a>
            </nav>
        </div>
    </nav>
    <div class="container-fluid">
        <div class="row">
            <div class="col-sm-3 col-md-2 sidebar">
                <ul class="nav nav-sidebar">
                    <li><a href="#dashboard">Dashboard</a></li>
                    <li><a href="#help"><i class="fa fa-question"></i></a></li>
                </ul>
            </div>
            <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
                <h2 class="sub-header" id="dashboard">Dashboard</h2>
                <br>
                <div class="row">
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd1_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd1_Name"><b>Name:</b> ...</li>
                                    <li id="etcd1_ID"><b>ID:</b> ...</li>
                                    <li id="etcd1_StartTime"><b>StartTime:</b> ...</li>
                                    <li id="etcd1_LeaderUptime"><b>LeaderUptime:</b> ...</li>
                                    <li id="etcd1_RecvAppendRequestCnt"><b>RecvAppendRequestCnt:</b> ...</li>
                                    <li id="etcd1_RecvingBandwidthRate"><b>RecvingBandwidthRate:</b> ...</li>
                                    <li id="etcd1_SendAppendRequestCnt"><b>SendAppendRequestCnt:</b> ...</li>
                                    <li id="etcd1_SendingBandwidthRate"><b>SendingBandwidthRate:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd1_state_button" class="btn btn-secondary"><i id="etcd1_state" class="fa fa-pause-circle"></i></button>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd2_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd2_Name"><b>Name:</b> ...</li>
                                    <li id="etcd2_ID"><b>ID:</b> ...</li>
                                    <li id="etcd2_StartTime"><b>StartTime:</b> ...</li>
                                    <li id="etcd2_LeaderUptime"><b>LeaderUptime:</b> ...</li>
                                    <li id="etcd2_RecvAppendRequestCnt"><b>RecvAppendRequestCnt:</b> ...</li>
                                    <li id="etcd2_RecvingBandwidthRate"><b>RecvingBandwidthRate:</b> ...</li>
                                    <li id="etcd2_SendAppendRequestCnt"><b>SendAppendRequestCnt:</b> ...</li>
                                    <li id="etcd2_SendingBandwidthRate"><b>SendingBandwidthRate:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd2_state_button" class="btn btn-secondary"><i id="etcd2_state" class="fa fa-pause-circle"></i></button>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd3_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd3_ID"><b>ID:</b> ...</li>
                                    <li id="etcd3_Name"><b>Name:</b> ...</li>
                                    <li id="etcd3_StartTime"><b>StartTime:</b> ...</li>
                                    <li id="etcd3_LeaderUptime"><b>LeaderUptime:</b> ...</li>
                                    <li id="etcd3_RecvAppendRequestCnt"><b>RecvAppendRequestCnt:</b> ...</li>
                                    <li id="etcd3_RecvingBandwidthRate"><b>RecvingBandwidthRate:</b> ...</li>
                                    <li id="etcd3_SendAppendRequestCnt"><b>SendAppendRequestCnt:</b> ...</li>
                                    <li id="etcd3_SendingBandwidthRate"><b>SendingBandwidthRate:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd3_state_button" class="btn btn-secondary"><i id="etcd3_state" class="fa fa-pause-circle"></i></button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="row">
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd4_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd4_Name"><b>Name:</b> ...</li>
                                    <li id="etcd4_ID"><b>ID:</b> ...</li>
                                    <li id="etcd4_StartTime"><b>StartTime:</b> ...</li>
                                    <li id="etcd4_LeaderUptime"><b>LeaderUptime:</b> ...</li>
                                    <li id="etcd4_RecvAppendRequestCnt"><b>RecvAppendRequestCnt:</b> ...</li>
                                    <li id="etcd4_RecvingBandwidthRate"><b>RecvingBandwidthRate:</b> ...</li>
                                    <li id="etcd4_SendAppendRequestCnt"><b>SendAppendRequestCnt:</b> ...</li>
                                    <li id="etcd4_SendingBandwidthRate"><b>SendingBandwidthRate:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd4_state_button" class="btn btn-secondary"><i id="etcd4_state" class="fa fa-pause-circle"></i></button>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd5_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd5_Name"><b>Name:</b> ...</li>
                                    <li id="etcd5_ID"><b>ID:</b> ...</li>
                                    <li id="etcd5_StartTime"><b>StartTime:</b> ...</li>
                                    <li id="etcd5_LeaderUptime"><b>LeaderUptime:</b> ...</li>
                                    <li id="etcd5_RecvAppendRequestCnt"><b>RecvAppendRequestCnt:</b> ...</li>
                                    <li id="etcd5_RecvingBandwidthRate"><b>RecvingBandwidthRate:</b> ...</li>
                                    <li id="etcd5_SendAppendRequestCnt"><b>SendAppendRequestCnt:</b> ...</li>
                                    <li id="etcd5_SendingBandwidthRate"><b>SendingBandwidthRate:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd5_state_button" class="btn btn-secondary"><i id="etcd5_state" class="fa fa-pause-circle"></i></button>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <div class="input-group">
                                    <form method="POST" id="endpoint">
                                        <textarea name="endpoint_name" type="text" class="form-control" placeholder="New-line separated endpoints (max 5)" rows="8"></textarea>
                                    </form>
                                    <span class="input-group-btn">
                                <button class="btn btn-secondary" form="endpoint" type="submit">Submit</button>
                            </span>
                                </div>
                            </p>
                            <div class="wrapper">
                                <button id="endpoint_state_button" class="btn btn-secondary"><i id="endpoint_state" class="fa fa-pause-circle"></i></button>
                                <span class="inline" id="endpoint_result">&nbsp;&nbsp;Not yet...</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="portlet">
                    <div class="portlet-body">
                        <div class="row">
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_StorageKeysTotal" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_StorageBytes" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_WalFsyncSecondsSum" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_GcSecondsSum" style="margin: 0 auto"></div>
                                </div>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_MemstatsAllocBytes" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_MemstatsHeapAllocBytes" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_MemstatsMallocsTotal" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-3">
                                <div class="card card-block">
                                    <div id="chart_ProcessCPUSeconds" style="margin: 0 auto"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-md-3">
                            <div class="card card-block">
                                <div id="chart_Goroutines" style="margin: 0 auto"></div>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="card card-block">
                                <div id="chart_StorageWatcherTotal" style="margin: 0 auto"></div>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="card card-block">
                                <div id="chart_StorageWatchStreamTotal" style="margin: 0 auto"></div>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="card card-block">
                                <div id="chart_StorageSlowWatcherTotal" style="margin: 0 auto"></div>
                            </div>
                        </div>
                    </div>
                </div>
                <br>
                <br>
                <h2 class="sub-header" id="help">Help</h2>
                <div class="row">
                    <div class="col-md-9">
                        <style type="text/css">
                        .tg {
                            border-collapse: collapse;
                            border-spacing: 0;
                            border-color: #ccc;
                        }
                        
                        .tg td {
                            font-family: Arial, sans-serif;
                            font-size: 14px;
                            padding: 10px 5px;
                            border-style: solid;
                            border-width: 1px;
                            overflow: hidden;
                            word-break: normal;
                            border-color: #ccc;
                            color: #333;
                            background-color: #fff;
                        }
                        
                        .tg th {
                            font-family: Arial, sans-serif;
                            font-size: 14px;
                            font-weight: normal;
                            padding: 10px 5px;
                            border-style: solid;
                            border-width: 1px;
                            overflow: hidden;
                            word-break: normal;
                            border-color: #ccc;
                            color: #333;
                            background-color: #f0f0f0;
                        }
                        
                        .tg .tg-1 {
                            text-align: center;
                            margin: auto;
                            vertical-align: center;
                        }
                        
                        .tg .tg-2 {
                            padding-left: 10px;
                            padding-right: 7px;
                            vertical-align: center;
                        }
                        </style>
                        <table class="tg">
                            <tr>
                                <th class="tg-1">Button</th>
                                <th class="tg-1">Description</th>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-secondary btn-sm"><i class="fa fa-pause-circle"></i></button>
                                </td>
                                <td class="tg-2">Node state / Endpoints status before start</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-danger btn-sm"><i class="fa fa-frown-o"></i></button>
                                </td>
                                <td class="tg-2">Node is not reachable (killed or error-ed)</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-warning btn-sm"><i class="fa fa-child"></i></button>
                                </td>
                                <td class="tg-2">Node is the follower (StateFollower)</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-info btn-sm"><i class="fa fa-star-half-o"></i></button>
                                </td>
                                <td class="tg-2">Node is the candidate (StateCandidate)</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-primary btn-sm"><i class="fa fa-fort-awesome"></i></button>
                                </td>
                                <td class="tg-2">Node is the leader (StateLeader)</td>
                            </tr>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="https://code.jquery.com/ui/1.12.0-beta.1/jquery-ui.js"></script>
    <link rel="stylesheet" type="text/css" href="https://code.jquery.com/ui/1.11.4/themes/smoothness/jquery-ui.css">
</body>

</html>
`
