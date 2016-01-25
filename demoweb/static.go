package demoweb

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
        function split(val) {
            return val.split(/,\s*/);
        }

        function extractLast(term) {
            return split(term).pop();
        }

        function appendLog(msg) {
            msg = $("<div />").html(msg);
            $("#log_output").append(msg);
            $("#log_box").scrollTop($("#log_box")[0].scrollHeight);
        }

        function runWithTimeout(callback, delay, times) {
            var internalCallback = function(t) {
                return function() {
                    if (--t > 0) {
                        window.setTimeout(internalCallback, delay);
                        callback();
                    }
                }
            }(times);
            window.setTimeout(internalCallback, delay);
        };

        var conn;
        var wsURL = "ws://localhost:8000/ws";
        if (window["WebSocket"]) {
            conn = new WebSocket(wsURL);
            conn.onopen = function() {
                console.log("connected to", wsURL);
            }
            conn.onclose = function(ev) {
                appendLog($("<div><b>connection closed</b></div>"));
            }
            conn.onmessage = function(ev) {
                appendLog(ev.data);
            }
            conn.onerror = function(ev) {
                appendLog("ERROR: " + ev.data);
            }
        } else {
            appendLog($("<div><b>browser does not support WebSocket</b></div>"))
        }

        var receiveStream = function() {
            $.ajax({
                url: '/stream',
                dataType: "json",
                success: function(dataObj) {
                    if (dataObj.Size > 0) {
                        // conn.send(dataObj.Logs);
                        appendLog(dataObj.Logs);
                    }
                }
            });
        }

        var receiveStats = function() {
            $.ajax({
                type: "GET",
                url: '/stats',
                async: true,
                dataType: "json",
                success: function(dataObj) {
                    document.getElementById('etcd1_Name').innerHTML = "Name: <b>" + dataObj.Etcd1Name + "</b>";
                    document.getElementById('etcd1_Endpoint').innerHTML = "Endpoint: <b>" + dataObj.Etcd1Endpoint + "</b>";
                    document.getElementById('etcd1_ID').innerHTML = "ID: <b>" + dataObj.Etcd1ID + "</b>";
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
                }
            });
        }

        Highcharts.setOptions({
            global: {
                useUTC: false
            }
        });

        $('#chart_StorageKeysTotal').highcharts({
            chart: {
                type: 'spline',
                animation: Highcharts.svg, // don't animate in old IE
                marginRight: 10,
                events: {
                    load: function() {
                        // var series1 = this.series[0];
                        // var series2 = this.series[1];
                        // var series3 = this.series[2];
                        // setInterval(function() {
                        //     $.ajax({
                        //         type: "GET",
                        //         url: '/metrics',
                        //         async: true,
                        //         dataType: "json",
                        //         success: function(dataObj) {
                        //             var x = (new Date()).getTime();
                        //             series1.addPoint({
                        //                 x: x,
                        //                 y: dataObj.Etcd1StorageKeysTotal
                        //             }, true, true);
                        //             series2.addPoint({
                        //                 x: x,
                        //                 y: dataObj.Etcd2StorageKeysTotal
                        //             }, true, true);
                        //             series3.addPoint({
                        //                 x: x,
                        //                 y: dataObj.Etcd3StorageKeysTotal
                        //             }, true, true);
                        //         }
                        //     });
                        // }, 5000);

                        var series1 = this.series[0];
                        var series2 = this.series[1];
                        var series3 = this.series[2];
                        var receiveMetrics = function() {
                            $.ajax({
                                type: "GET",
                                url: '/metrics',
                                async: true,
                                dataType: "json",
                                success: function(dataObj) {
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
                                }
                            });
                        };
                        runWithTimeout(receiveMetrics, 1000, 1);
                        runWithTimeout(receiveMetrics, 10000, 1);
                        runWithTimeout(receiveMetrics, 5000, 10000);
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
                            x: time + i * 5000,
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
                            x: time + i * 5000,
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
                            x: time + i * 5000,
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
                        // var series1 = this.series[0];
                        // var series2 = this.series[1];
                        // var series3 = this.series[2];
                        // setInterval(function() {
                        //     $.ajax({
                        //         type: "GET",
                        //         url: '/metrics',
                        //         async: true,
                        //         dataType: "json",
                        //         success: function(dataObj) {
                        //             var x = (new Date()).getTime();
                        //             series1.addPoint({
                        //                 x: x,
                        //                 y: dataObj.Etcd1StorageBytes
                        //             }, true, true);
                        //             series2.addPoint({
                        //                 x: x,
                        //                 y: dataObj.Etcd2StorageBytes
                        //             }, true, true);
                        //             series3.addPoint({
                        //                 x: x,
                        //                 y: dataObj.Etcd3StorageBytes
                        //             }, true, true);
                        //         }
                        //     });
                        // }, 5000);

                        var series1 = this.series[0];
                        var series2 = this.series[1];
                        var series3 = this.series[2];
                        var receiveMetrics = function() {
                            $.ajax({
                                type: "GET",
                                url: '/metrics',
                                async: true,
                                dataType: "json",
                                success: function(dataObj) {
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
                                }
                            });
                        };
                        runWithTimeout(receiveMetrics, 100, 1);
                        runWithTimeout(receiveMetrics, 10000, 1);
                        runWithTimeout(receiveMetrics, 5000, 10000);
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
                            x: time + i * 5000,
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
                            x: time + i * 5000,
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
                            x: time + i * 5000,
                            y: 0,
                            bytes: "0 bytes"
                        });
                    }
                    return data;
                }())
            }]
        });

        // var eventStream = new EventSource('/stream')
        // eventStream.onopen = function(ev) {
        //     console.log("connected to", ev.target.url);
        // };
        // eventStream.onmessage = function(ev) {
        //     conn.send(ev.data);
        //     // appendLog(ev.data);
        // };
        // eventStream.addEventListener('stream_log', function(ev) {
        //     // conn.send(ev.data);
        //     appendLog(ev.data);
        // });

        // setInterval(function() {
        //     $.ajax({
        //         url: '/stream',
        //         success: function(ev) {
        //             conn.send(ev);
        //             // appendLog(ev);
        //         }
        //     });
        // }, 1000);

        // setInterval(function() {
        //     $.ajax({
        //         url: '/stream',
        //         dataType: "json",
        //         success: function(ev) {
        //             if (ev.Size > 0) {
        //                 conn.send(ev.Value);
        //                 // appendLog(ev.Value); 
        //             }
        //         }
        //     });
        // }, 1000);

        $('#start_cluster').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/start_cluster",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#start_cluster').click(function(e) {
            // TODO: dynamically change timeout
            runWithTimeout(receiveStream, 100, 1);
            runWithTimeout(receiveStream, 1000, 1);
            runWithTimeout(receiveStream, 7000, 24);
            runWithTimeout(receiveStream, 10000, 300);
            runWithTimeout(receiveStream, 30000, 500);
            runWithTimeout(receiveStream, 60000, 700);
            runWithTimeout(receiveStream, 100000, 10000);
        });
        $('#start_cluster').click(function(e) {
            runWithTimeout(receiveStats, 100, 1);
            runWithTimeout(receiveStats, 2000, 3);
            runWithTimeout(receiveStats, 10000, 1000);
        });

        $('#start_stress').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/start_stress",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#start_stress').click(function(e) {
            runWithTimeout(receiveStats, 100, 1);
            runWithTimeout(receiveStats, 1000, 1);
            runWithTimeout(receiveStats, 300, 2);
        });
        $('#start_stress').click(function(e) {
            runWithTimeout(receiveStream, 1500, 1);
        });

        $("#list_ctl")
            // don't navigate away from the field on tab when selecting an item
            .on("click", function(event, ui) {
                if (event.keyCode === $.ui.keyCode.TAB &&
                    $(this).autocomplete("instance").menu.active) {
                    event.preventDefault();
                }
            })
            .autocomplete({
                minLength: 0,
                source: function(request, response) {
                    $.ajax({
                        type: "GET",
                        url: "/list_ctl",
                        dataType: "json",
                        data: {
                            q: request.term
                        },
                        success: function(dataObj) {
                                response($.ui.autocomplete.filter(dataObj.Values, extractLast(request.term)));
                            }
                            // delegate back to autocomplete, but extract the last term
                            // response($.ui.autocomplete.filter(availableVerticals, extractLast(request.term)));
                    });
                },
                focus: function() {
                    // prevent value inserted on focus
                    return false;
                }
            });

        $('#ctl').submit(function(e) {
            e.preventDefault();
            $.ajax({
                type: "POST",
                url: "/ctl",
                data: $(this).serialize(),
                success: function(dataObj) {
                    $.ajax({
                        type: "GET",
                        url: "/ctl",
                        async: true,
                        dataType: "html",
                        success: function(dataObj) {
                            conn.send(dataObj);
                            // appendLog(dataObj);
                        }
                    });
                }
            });
        });

        $('#kill_1').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/kill_1",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#kill_1').click(function(e) {
            runWithTimeout(receiveStats, 100, 2);
            runWithTimeout(receiveStats, 500, 1);
        });
        $('#kill_1').click(function(e) {
            runWithTimeout(receiveStream, 700, 1);
        });

        $('#restart_1').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/restart_1",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#restart_1').click(function(e) {
            runWithTimeout(receiveStats, 100, 2);
            runWithTimeout(receiveStats, 500, 1);
        });
        $('#restart_1').click(function(e) {
            runWithTimeout(receiveStream, 700, 1);
        });

        $('#kill_2').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/kill_2",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#kill_2').click(function(e) {
            runWithTimeout(receiveStats, 100, 2);
            runWithTimeout(receiveStats, 500, 1);
        });
        $('#kill_2').click(function(e) {
            runWithTimeout(receiveStream, 700, 1);
        });

        $('#restart_2').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/restart_2",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#restart_2').click(function(e) {
            runWithTimeout(receiveStats, 100, 2);
            runWithTimeout(receiveStats, 500, 1);
        });
        $('#restart_2').click(function(e) {
            runWithTimeout(receiveStream, 700, 1);
        });

        $('#kill_3').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/kill_3",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#kill_3').click(function(e) {
            runWithTimeout(receiveStats, 100, 2);
            runWithTimeout(receiveStats, 500, 1);
        });
        $('#kill_3').click(function(e) {
            runWithTimeout(receiveStream, 700, 1);
        });

        $('#restart_3').click(function(e) {
            e.preventDefault();
            $.ajax({
                type: "GET",
                url: "/restart_3",
                async: true,
                dataType: "html",
                success: function(dataObj) {
                    // conn.send(dataObj);
                    appendLog(dataObj);
                }
            });
        });
        $('#restart_3').click(function(e) {
            runWithTimeout(receiveStats, 100, 2);
            runWithTimeout(receiveStats, 500, 1);
        });
        $('#restart_3').click(function(e) {
            runWithTimeout(receiveStream, 700, 1);
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
            width: 100px;
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
    
    #log_box {
        font-family: "Lucida Console", Courier, monospace;
        font-size: 9px;
        font-style: normal;
        font-variant: normal;
    }
    </style>
    <title>runetcd demo-web</title>
</head>

<body>
    <nav class="navbar navbar-dark navbar-fixed-top bg-inverse">
        <button type="button" class="navbar-toggler hidden-sm-up" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
        </button>
        <a class="navbar-brand" href="https://github.com/gophergala2016/runetcd">Run etcd</a>
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
                    <li><a href="#demo">Demo</a></li>
                    <li><a href="#log">Log</a></li>
                    <li><a href="#help"><i class="fa fa-question"></i></a></li>
                </ul>
            </div>
            <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
                <h2 class="sub-header" id="demo">Demo</h2>
                <br>
                <div class="row">
                    <div class="col-md-3">
                        <div class="wrapper">
                            <input type="submit" class="btn btn-secondary btn-sm" id="start_cluster" value="Start">
                            <input type="submit" class="btn btn-secondary btn-sm" id="start_stress" value="Stress">
                        </div>
                    </div>
                    <div class="col-md-7">
                        <div class="input-group">
                            <form method="POST" id="ctl">
                                <input id="list_ctl" name="ctl_name" type="text" class="form-control" placeholder="etcdctlv3 ...">
                            </form>
                            <span class="input-group-btn">
                            <button class="btn btn-default" form="ctl" type="submit">Submit</button>
                            </span>
                        </div>
                    </div>
                </div>
                <br>
                <div class="row">
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd1_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd1_Name"><b>Name:</b> ...</li>
                                    <li id="etcd1_ID"><b>ID:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd1_state_button" class="btn btn-secondary-outline"><i id="etcd1_state" class="fa fa-pause-circle"></i></button>
                                <button class="btn btn-danger" id="kill_1"><i class="fa fa-bomb"></i></button>
                                <button class="btn btn-success" id="restart_1"><i class="fa fa-ambulance"></i></button>
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
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd2_state_button" class="btn btn-secondary-outline"><i id="etcd2_state" class="fa fa-pause-circle"></i></button>
                                <button class="btn btn-danger" id="kill_2"><i class="fa fa-bomb"></i></button>
                                <button class="btn btn-success" id="restart_2"><i class="fa fa-ambulance"></i></button>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card card-block">
                            <p class="card-text">
                                <ul>
                                    <li id="etcd3_Endpoint"><b>Endpoint:</b> ...</li>
                                    <li id="etcd3_Name"><b>Name:</b> ...</li>
                                    <li id="etcd3_ID"><b>ID:</b> ...</li>
                                </ul>
                            </p>
                            <div class="wrapper">
                                <button id="etcd3_state_button" class="btn btn-secondary-outline"><i id="etcd3_state" class="fa fa-pause-circle"></i></button>
                                <button class="btn btn-danger" id="kill_3"><i class="fa fa-bomb"></i></button>
                                <button class="btn btn-success" id="restart_3"><i class="fa fa-ambulance"></i></button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="portlet">
                    <div class="portlet-body">
                        <div class="row">
                            <div class="col-md-6">
                                <div class="card card-block">
                                    <div id="chart_StorageKeysTotal" style="margin: 0 auto"></div>
                                </div>
                            </div>
                            <div class="col-md-6">
                                <div class="card card-block">
                                    <div id="chart_StorageBytes" style="margin: 0 auto"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <br>
                <h2 class="sub-header" id="log">Log</h2>
                <div id="log_box" style="padding-top: 10px; padding-left: 5px; border:1px solid black; height:500px; overflow: scroll;">
                    <div id="log_output">
                    </div>
                </div>
                <br>
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
                                    <button type="button" class="btn btn-secondary btn-sm">Start</button>
                                </td>
                                <td class="tg-2">Start the 3-member cluster</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-secondary btn-sm">Stress</button>
                                </td>
                                <td class="tg-2">Stress the cluster with random key-values</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-default btn-sm">Submit</button>
                                </td>
                                <td class="tg-2">Submit your own etcdctl commands</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-secondary-outline btn-sm"><i class="fa fa-pause-circle"></i></button>
                                </td>
                                <td class="tg-2">Node state before start</td>
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
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-danger btn-sm"><i class="fa fa-bomb"></i></button>
                                </td>
                                <td class="tg-2">Kill the node</td>
                            </tr>
                            <tr>
                                <td class="tg-1">
                                    <button type="button" class="btn btn-success btn-sm"><i class="fa fa-ambulance"></i></button>
                                </td>
                                <td class="tg-2">Recover the node</td>
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
