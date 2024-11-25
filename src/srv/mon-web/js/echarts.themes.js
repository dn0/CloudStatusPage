(function () {
  const darkMode = window?.matchMedia?.('(prefers-color-scheme:dark)')?.matches;

  window.activeCharts = new Map();

  window.echartTheme = function () {
    if (darkMode) {
      return 'dark2';
    }
    return 'light2';
  };

  window.debounce = function (callback, time) {
    var timeout;
    return function() {
      var context = this;
      var args = arguments;
      if (timeout) {
        clearTimeout(timeout);
      }
      timeout = setTimeout(function() {
        timeout = null;
        callback.apply(context, args);
      }, time);
    }
  };

  window.echartMapClickEventListener = function(params) {
    // Because the world wants to be complicated
    // (missing feature in htmx where there is no ajax api that would also update the browser history)
    if (params.value[3]) {
      const fakeLink = document.createElement("a");
      fakeLink.setAttribute('href', params.value[3]);
      htmx.process(fakeLink);
      fakeLink.click();
    }
  };

  window.echartMapTooltipFormatter = function(params) {
    let message = '<div class="chart-tooltip-wrapper">';
    message += '<b>'+ params.name +'</b>';
    message += '<br/>' + params.value[2];
    message += '</div>';
    return message;
  }

  window.echartLineTooltipFormatterFactory = function(unit) {
    return function(params) {
      if (!(params instanceof Array && params.length)) {
        return null;
      }
      let empty = true;
      let message = '<div class="chart-tooltip-wrapper">';
      message += '<span class="chart-tooltip-title">'+ params[0].axisValueLabel +'</span>';
      params.forEach(param => {
        if (param.value[1] === null) {
          return;
        }
        empty = false;
        message += '<br/>' + param.marker + param.seriesName + '<span class="chart-tooltip-value">'+ param.value[1] +' '+ unit +'</span>';
        if (param.value.length === 3) {
          message += '<pre class="chart-tooltip-code">'+ param.value[2] +'</pre>';
        }
      });
      if (empty) {
        return null;
      }
      message += '</div>';
      return message;
    }
  };

})();

/*
 * Theme: light2
 */
(function (root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['exports', 'echarts'], factory);
  } else if (typeof exports === 'object' && typeof exports.nodeName !== 'string') {
    // CommonJS
    factory(exports, require('echarts'));
  } else {
    // Browser globals
    factory({}, root.echarts);
  }
}(this, function (exports, echarts) {
    let log = function (msg) {
      if (typeof console !== 'undefined') {
        console && console.error && console.error(msg);
      }
    };
    if (!echarts) {
      log('ECharts is not Loaded');
      return;
    }

    const c50 = "#f9fafb",
      c100 = "#f3f4f6",
      c200 = "#e5e7eb",
      c300 = "#d1d5db",
      c400 = "#9ca3af",
      c500 = "#6b7280",
      c600 = "#4b5563",
      c700 = "#374151",
      c800 = "#1f2937",
      c900 = "#111827",
      c950 = "#030712",
      c1000 = "#0ea5e9";

    echarts.registerTheme('light2', {
      "color": [
        "#4992ff",
        "#7cffb2",
        "#fddd60",
        "#ff6e76",
        "#58d9f9",
        "#05c091",
        "#ff8a45",
        "#8d48e3",
        "#dd79ff"
      ],
      "backgroundColor": "white",
      "title": {
        "top": "8px",
        "left": "8px",
        "textStyle": {
          "color": c700,
          "fontSize": 15,
          "fontFamily": "inter, ui-sans-serif, sans-serif"
        },
        "subtextStyle": {
          "color": c500
        }
      },
      "tooltip": {
        "show": true,
        "enterable": true,
        "backgroundColor": c50,
        "textStyle": {
          "color": c700,
          "fontSize": 13,
          "fontFamily": "inter, ui-sans-serif, sans-serif"
        },
        "emphasis": {
          "focus": true
        },
        "axisPointer": {
          "show": true,
          "type": "line",
          "lineStyle": {
            "color": c200,
            "width": "1"
          },
          "crossStyle": {
            "color": c200,
            "width": "1"
          }
        }
      },
      "legend": {
        "show": true,
        "type": "plain",
        "top": "28px",
        "textStyle": {
          "color": c700,
          "fontSize": 12
        },
        "emphasis": {
          "focus": true
        }
      },
      "toolbox": {
        "show": true,
        "top": "8px",
        "right": "8px",
        "orient": "horizontal",
        "left": "right",
        "iconStyle": {
          "borderColor": c600
        },
        "emphasis": {
          "textStyle": {
            "color": c800,
            "borderColor": c700
          },
          "iconStyle": {
            "borderColor": c700
          }
        }
      },
      "dataZoom": {
        "type": "slider",
        "backgroundColor": "transparent",
        "borderColor": c200,
        "textStyle": {
          "color": c950,
          "fontWeight": 600
        },
        "fillerColor": "rgba(114,204,255,0.2)",
        "dataBackground": {
          "lineStyle": {
            "color": c500,
          },
          "areaStyle": {
            "color": c600,
            "opacity": 0.3
          }
        },
        "selectedDataBackground": {
          "lineStyle": {
            "color": c1000
          },
          "areaStyle": {
            "color": c1000,
            "opacity": 0.5
          }
        },
      },
      "markLine": {
        "label": {
          "color": c700,
          "lineHeight": 16
        },
        "lineStyle": {
          "opacity": 0.9
        }
      },
      "markArea": {
        "itemStyle": {
          "color": "#f87171",
          "opacity": 0.2
        }
      },
      "grid": {
        "top": "76px",
        "left": "60px",
        "right": "60px",
        "bottom": "70px",
        "show": true,
        "backgroundColor": c50,
        "borderWidth": 0,
        "borderColor": "transparent"
      },
      "valueAxis": {
        "axisLine": {
          "show": true,
          "lineStyle": {
            "color": c500,
            "width": 0.7
          }
        },
        "axisLabel": {
          "show": true,
          "color": c700
        },
        "splitLine": {
          "show": true,
          "lineStyle": {
            "color": [
              c200
            ]
          }
        }
      },
      "timeAxis": {
        "axisLine": {
          "show": true,
          "lineStyle": {
            "color": c500,
            "width": 0.7
          }
        },
        "axisLabel": {
          "show": true,
          "color": c700
        }
      },
      "effectScatter": {
        "zlevel": 5,
        "z": 5
      },
      "geo": {
        "itemStyle": {
          "areaColor": c400,
          "borderColor": c50,
          "borderWidth": 0.5
        },
        "label": {
          "color": c700
        }
      }
    });
  }));

/*
 * Theme: dark2
 */
(function (root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['exports', 'echarts'], factory);
  } else if (typeof exports === 'object' && typeof exports.nodeName !== 'string') {
    // CommonJS
    factory(exports, require('echarts'));
  } else {
    // Browser globals
    factory({}, root.echarts);
  }
}(this, function (exports, echarts) {
  let log = function (msg) {
    if (typeof console !== 'undefined') {
      console && console.error && console.error(msg);
    }
  };
  if (!echarts) {
    log('ECharts is not Loaded');
    return;
  }

  const c50 = "#f9fafb",
    c100 = "#f3f4f6",
    c200 = "#e5e7eb",
    c300 = "#d1d5db",
    c400 = "#9ca3af",
    c500 = "#6b7280",
    c600 = "#4b5563",
    c700 = "#374151",
    c800 = "#1f2937",
    c900 = "#111827",
    c950 = "#030712",
    c1000 = "#38bdf8";

  echarts.registerTheme('dark2', {
    "color": [
      "#4992ff",
      "#05c091",
      "#fddd60",
      "#ff6e76",
      "#58d9f9",
      "#7cffb2",
      "#ff8a45",
      "#8d48e3",
      "#dd79ff"
    ],
    "backgroundColor": "black",
    "title": {
      "top": "8px",
      "left": "8px",
      "textStyle": {
        "color": c200,
        "fontSize": 15,
        "fontFamily": "inter, ui-sans-serif, sans-serif"
      },
      "subtextStyle": {
        "color": c400
      }
    },
    "tooltip": {
      "show": true,
      "enterable": true,
      "backgroundColor": c950,
      "textStyle": {
        "color": c200,
        "fontSize": 13,
        "fontFamily": "inter, ui-sans-serif, sans-serif"
      },
      "emphasis": {
        "focus": true
      },
      "axisPointer": {
        "show": true,
        "type": "line",
        "lineStyle": {
          "color": c700,
          "width": "1"
        },
        "crossStyle": {
          "color": c700,
          "width": "1"
        }
      }
    },
    "legend": {
      "show": true,
      "type": "plain",
      "top": "28px",
      "textStyle": {
        "color": c200,
        "fontSize": 12
      },
      "emphasis": {
        "focus": true
      }
    },
    "toolbox": {
      "show": true,
      "top": "8px",
      "right": "8px",
      "orient": "horizontal",
      "left": "right",
      "iconStyle": {
        "borderColor": c300
      },
      "emphasis": {
        "textStyle": {
          "color": c100,
          "borderColor": c200
        },
        "iconStyle": {
          "borderColor": c200
        }
      }
    },
    "dataZoom": {
      "type": "slider",
      "backgroundColor": "transparent",
      "borderColor": c700,
      "textStyle": {
        "color": c50,
        "fontWeight": 600
      },
      "fillerColor": "rgba(114,204,255,0.2)",
      "dataBackground": {
        "lineStyle": {
          "color": c400,
        },
        "areaStyle": {
          "color": c300,
          "opacity": 0.8
        }
      },
      "selectedDataBackground": {
        "lineStyle": {
          "color": c1000
        },
        "areaStyle": {
          "color": c1000,
          "opacity": 0.5
        }
      },
    },
    "markLine": {
      "label": {
        "color": c200,
        "lineHeight": 16
      },
      "lineStyle": {
        "opacity": 0.9
      }
    },
    "markArea": {
      "itemStyle": {
        "color": "#f87171",
        "opacity": 0.2
      }
    },
    "grid": {
      "top": "76px",
      "left": "60px",
      "right": "60px",
      "bottom": "70px",
      "show": true,
      "backgroundColor": c950,
      "borderWidth": 0,
      "borderColor": "transparent"
    },
    "valueAxis": {
      "axisLine": {
        "show": true,
        "lineStyle": {
          "color": c400,
          "width": 0.7
        }
      },
      "axisLabel": {
        "show": true,
        "color": c200
      },
      "splitLine": {
        "show": true,
        "lineStyle": {
          "color": [
            c800
          ]
        }
      }
    },
    "timeAxis": {
      "axisLine": {
        "show": true,
        "lineStyle": {
          "color": c400,
          "width": 0.7
        }
      },
      "axisLabel": {
        "show": true,
        "color": c200
      }
    },
    "effectScatter": {
      "zlevel": 5,
      "z": 5
    },
    "geo": {
      "itemStyle": {
        "areaColor": c400,
        "borderColor": c950,
        "borderWidth": 0.5
      },
      "label": {
        "color": c200
      }
    }
  });
}));
