var __create = Object.create;
var __defProp = Object.defineProperty;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __markAsModule = (target) => __defProp(target, "__esModule", {value: true});
var __commonJS = (callback, module2) => () => {
  if (!module2) {
    module2 = {exports: {}};
    callback(module2.exports, module2);
  }
  return module2.exports;
};
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, {get: all[name], enumerable: true});
};
var __exportStar = (target, module2, desc) => {
  if (module2 && typeof module2 === "object" || typeof module2 === "function") {
    for (let key of __getOwnPropNames(module2))
      if (!__hasOwnProp.call(target, key) && key !== "default")
        __defProp(target, key, {get: () => module2[key], enumerable: !(desc = __getOwnPropDesc(module2, key)) || desc.enumerable});
  }
  return target;
};
var __toModule = (module2) => {
  return __exportStar(__markAsModule(__defProp(module2 != null ? __create(__getProtoOf(module2)) : {}, "default", module2 && module2.__esModule && "default" in module2 ? {get: () => module2.default, enumerable: true} : {value: module2, enumerable: true})), module2);
};

// node_modules/object-assign/index.js
var require_object_assign = __commonJS((exports2, module2) => {
  /*
  object-assign
  (c) Sindre Sorhus
  @license MIT
  */
  "use strict";
  var getOwnPropertySymbols = Object.getOwnPropertySymbols;
  var hasOwnProperty = Object.prototype.hasOwnProperty;
  var propIsEnumerable = Object.prototype.propertyIsEnumerable;
  function toObject(val) {
    if (val === null || val === void 0) {
      throw new TypeError("Object.assign cannot be called with null or undefined");
    }
    return Object(val);
  }
  function shouldUseNative() {
    try {
      if (!Object.assign) {
        return false;
      }
      var test1 = new String("abc");
      test1[5] = "de";
      if (Object.getOwnPropertyNames(test1)[0] === "5") {
        return false;
      }
      var test2 = {};
      for (var i = 0; i < 10; i++) {
        test2["_" + String.fromCharCode(i)] = i;
      }
      var order2 = Object.getOwnPropertyNames(test2).map(function(n) {
        return test2[n];
      });
      if (order2.join("") !== "0123456789") {
        return false;
      }
      var test3 = {};
      "abcdefghijklmnopqrst".split("").forEach(function(letter) {
        test3[letter] = letter;
      });
      if (Object.keys(Object.assign({}, test3)).join("") !== "abcdefghijklmnopqrst") {
        return false;
      }
      return true;
    } catch (err) {
      return false;
    }
  }
  module2.exports = shouldUseNative() ? Object.assign : function(target, source) {
    var from;
    var to = toObject(target);
    var symbols;
    for (var s = 1; s < arguments.length; s++) {
      from = Object(arguments[s]);
      for (var key in from) {
        if (hasOwnProperty.call(from, key)) {
          to[key] = from[key];
        }
      }
      if (getOwnPropertySymbols) {
        symbols = getOwnPropertySymbols(from);
        for (var i = 0; i < symbols.length; i++) {
          if (propIsEnumerable.call(from, symbols[i])) {
            to[symbols[i]] = from[symbols[i]];
          }
        }
      }
    }
    return to;
  };
});

// node_modules/react/cjs/react.production.min.js
var require_react_production_min = __commonJS((exports2) => {
  /** @license React v17.0.1
   * react.production.min.js
   *
   * Copyright (c) Facebook, Inc. and its affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   */
  "use strict";
  var l = require_object_assign();
  var n = 60103;
  var p = 60106;
  exports2.Fragment = 60107;
  exports2.StrictMode = 60108;
  exports2.Profiler = 60114;
  var q = 60109;
  var r = 60110;
  var t = 60112;
  exports2.Suspense = 60113;
  var u = 60115;
  var v = 60116;
  if (typeof Symbol === "function" && Symbol.for) {
    w = Symbol.for;
    n = w("react.element");
    p = w("react.portal");
    exports2.Fragment = w("react.fragment");
    exports2.StrictMode = w("react.strict_mode");
    exports2.Profiler = w("react.profiler");
    q = w("react.provider");
    r = w("react.context");
    t = w("react.forward_ref");
    exports2.Suspense = w("react.suspense");
    u = w("react.memo");
    v = w("react.lazy");
  }
  var w;
  var x = typeof Symbol === "function" && Symbol.iterator;
  function y(a) {
    if (a === null || typeof a !== "object")
      return null;
    a = x && a[x] || a["@@iterator"];
    return typeof a === "function" ? a : null;
  }
  function z(a) {
    for (var b = "https://reactjs.org/docs/error-decoder.html?invariant=" + a, c = 1; c < arguments.length; c++)
      b += "&args[]=" + encodeURIComponent(arguments[c]);
    return "Minified React error #" + a + "; visit " + b + " for the full message or use the non-minified dev environment for full errors and additional helpful warnings.";
  }
  var A = {isMounted: function() {
    return false;
  }, enqueueForceUpdate: function() {
  }, enqueueReplaceState: function() {
  }, enqueueSetState: function() {
  }};
  var B = {};
  function C(a, b, c) {
    this.props = a;
    this.context = b;
    this.refs = B;
    this.updater = c || A;
  }
  C.prototype.isReactComponent = {};
  C.prototype.setState = function(a, b) {
    if (typeof a !== "object" && typeof a !== "function" && a != null)
      throw Error(z(85));
    this.updater.enqueueSetState(this, a, b, "setState");
  };
  C.prototype.forceUpdate = function(a) {
    this.updater.enqueueForceUpdate(this, a, "forceUpdate");
  };
  function D() {
  }
  D.prototype = C.prototype;
  function E(a, b, c) {
    this.props = a;
    this.context = b;
    this.refs = B;
    this.updater = c || A;
  }
  var F = E.prototype = new D();
  F.constructor = E;
  l(F, C.prototype);
  F.isPureReactComponent = true;
  var G = {current: null};
  var H = Object.prototype.hasOwnProperty;
  var I = {key: true, ref: true, __self: true, __source: true};
  function J(a, b, c) {
    var e, d = {}, k = null, h = null;
    if (b != null)
      for (e in b.ref !== void 0 && (h = b.ref), b.key !== void 0 && (k = "" + b.key), b)
        H.call(b, e) && !I.hasOwnProperty(e) && (d[e] = b[e]);
    var g = arguments.length - 2;
    if (g === 1)
      d.children = c;
    else if (1 < g) {
      for (var f = Array(g), m = 0; m < g; m++)
        f[m] = arguments[m + 2];
      d.children = f;
    }
    if (a && a.defaultProps)
      for (e in g = a.defaultProps, g)
        d[e] === void 0 && (d[e] = g[e]);
    return {$$typeof: n, type: a, key: k, ref: h, props: d, _owner: G.current};
  }
  function K(a, b) {
    return {$$typeof: n, type: a.type, key: b, ref: a.ref, props: a.props, _owner: a._owner};
  }
  function L(a) {
    return typeof a === "object" && a !== null && a.$$typeof === n;
  }
  function escape(a) {
    var b = {"=": "=0", ":": "=2"};
    return "$" + a.replace(/[=:]/g, function(a2) {
      return b[a2];
    });
  }
  var M = /\/+/g;
  function N(a, b) {
    return typeof a === "object" && a !== null && a.key != null ? escape("" + a.key) : b.toString(36);
  }
  function O(a, b, c, e, d) {
    var k = typeof a;
    if (k === "undefined" || k === "boolean")
      a = null;
    var h = false;
    if (a === null)
      h = true;
    else
      switch (k) {
        case "string":
        case "number":
          h = true;
          break;
        case "object":
          switch (a.$$typeof) {
            case n:
            case p:
              h = true;
          }
      }
    if (h)
      return h = a, d = d(h), a = e === "" ? "." + N(h, 0) : e, Array.isArray(d) ? (c = "", a != null && (c = a.replace(M, "$&/") + "/"), O(d, b, c, "", function(a2) {
        return a2;
      })) : d != null && (L(d) && (d = K(d, c + (!d.key || h && h.key === d.key ? "" : ("" + d.key).replace(M, "$&/") + "/") + a)), b.push(d)), 1;
    h = 0;
    e = e === "" ? "." : e + ":";
    if (Array.isArray(a))
      for (var g = 0; g < a.length; g++) {
        k = a[g];
        var f = e + N(k, g);
        h += O(k, b, c, f, d);
      }
    else if (f = y(a), typeof f === "function")
      for (a = f.call(a), g = 0; !(k = a.next()).done; )
        k = k.value, f = e + N(k, g++), h += O(k, b, c, f, d);
    else if (k === "object")
      throw b = "" + a, Error(z(31, b === "[object Object]" ? "object with keys {" + Object.keys(a).join(", ") + "}" : b));
    return h;
  }
  function P(a, b, c) {
    if (a == null)
      return a;
    var e = [], d = 0;
    O(a, e, "", "", function(a2) {
      return b.call(c, a2, d++);
    });
    return e;
  }
  function Q(a) {
    if (a._status === -1) {
      var b = a._result;
      b = b();
      a._status = 0;
      a._result = b;
      b.then(function(b2) {
        a._status === 0 && (b2 = b2.default, a._status = 1, a._result = b2);
      }, function(b2) {
        a._status === 0 && (a._status = 2, a._result = b2);
      });
    }
    if (a._status === 1)
      return a._result;
    throw a._result;
  }
  var R = {current: null};
  function S() {
    var a = R.current;
    if (a === null)
      throw Error(z(321));
    return a;
  }
  var T = {ReactCurrentDispatcher: R, ReactCurrentBatchConfig: {transition: 0}, ReactCurrentOwner: G, IsSomeRendererActing: {current: false}, assign: l};
  exports2.Children = {map: P, forEach: function(a, b, c) {
    P(a, function() {
      b.apply(this, arguments);
    }, c);
  }, count: function(a) {
    var b = 0;
    P(a, function() {
      b++;
    });
    return b;
  }, toArray: function(a) {
    return P(a, function(a2) {
      return a2;
    }) || [];
  }, only: function(a) {
    if (!L(a))
      throw Error(z(143));
    return a;
  }};
  exports2.Component = C;
  exports2.PureComponent = E;
  exports2.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED = T;
  exports2.cloneElement = function(a, b, c) {
    if (a === null || a === void 0)
      throw Error(z(267, a));
    var e = l({}, a.props), d = a.key, k = a.ref, h = a._owner;
    if (b != null) {
      b.ref !== void 0 && (k = b.ref, h = G.current);
      b.key !== void 0 && (d = "" + b.key);
      if (a.type && a.type.defaultProps)
        var g = a.type.defaultProps;
      for (f in b)
        H.call(b, f) && !I.hasOwnProperty(f) && (e[f] = b[f] === void 0 && g !== void 0 ? g[f] : b[f]);
    }
    var f = arguments.length - 2;
    if (f === 1)
      e.children = c;
    else if (1 < f) {
      g = Array(f);
      for (var m = 0; m < f; m++)
        g[m] = arguments[m + 2];
      e.children = g;
    }
    return {
      $$typeof: n,
      type: a.type,
      key: d,
      ref: k,
      props: e,
      _owner: h
    };
  };
  exports2.createContext = function(a, b) {
    b === void 0 && (b = null);
    a = {$$typeof: r, _calculateChangedBits: b, _currentValue: a, _currentValue2: a, _threadCount: 0, Provider: null, Consumer: null};
    a.Provider = {$$typeof: q, _context: a};
    return a.Consumer = a;
  };
  exports2.createElement = J;
  exports2.createFactory = function(a) {
    var b = J.bind(null, a);
    b.type = a;
    return b;
  };
  exports2.createRef = function() {
    return {current: null};
  };
  exports2.forwardRef = function(a) {
    return {$$typeof: t, render: a};
  };
  exports2.isValidElement = L;
  exports2.lazy = function(a) {
    return {$$typeof: v, _payload: {_status: -1, _result: a}, _init: Q};
  };
  exports2.memo = function(a, b) {
    return {$$typeof: u, type: a, compare: b === void 0 ? null : b};
  };
  exports2.useCallback = function(a, b) {
    return S().useCallback(a, b);
  };
  exports2.useContext = function(a, b) {
    return S().useContext(a, b);
  };
  exports2.useDebugValue = function() {
  };
  exports2.useEffect = function(a, b) {
    return S().useEffect(a, b);
  };
  exports2.useImperativeHandle = function(a, b, c) {
    return S().useImperativeHandle(a, b, c);
  };
  exports2.useLayoutEffect = function(a, b) {
    return S().useLayoutEffect(a, b);
  };
  exports2.useMemo = function(a, b) {
    return S().useMemo(a, b);
  };
  exports2.useReducer = function(a, b, c) {
    return S().useReducer(a, b, c);
  };
  exports2.useRef = function(a) {
    return S().useRef(a);
  };
  exports2.useState = function(a) {
    return S().useState(a);
  };
  exports2.version = "17.0.1";
});

// node_modules/react/index.js
var require_react = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_react_production_min();
  } else {
    module2.exports = null;
  }
});

// node_modules/react-dom/cjs/react-dom-server.node.production.min.js
var require_react_dom_server_node_production_min = __commonJS((exports2) => {
  /** @license React v17.0.1
   * react-dom-server.node.production.min.js
   *
   * Copyright (c) Facebook, Inc. and its affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   */
  "use strict";
  var l = require_object_assign();
  var n = require_react();
  var aa = require("stream");
  function p(a) {
    for (var b = "https://reactjs.org/docs/error-decoder.html?invariant=" + a, c = 1; c < arguments.length; c++)
      b += "&args[]=" + encodeURIComponent(arguments[c]);
    return "Minified React error #" + a + "; visit " + b + " for the full message or use the non-minified dev environment for full errors and additional helpful warnings.";
  }
  var q = 60106;
  var r = 60107;
  var u = 60108;
  var z = 60114;
  var B = 60109;
  var ba = 60110;
  var ca = 60112;
  var D = 60113;
  var da = 60120;
  var ea = 60115;
  var fa = 60116;
  var ha = 60121;
  var ia = 60117;
  var ja = 60119;
  var ka = 60129;
  var la = 60131;
  if (typeof Symbol === "function" && Symbol.for) {
    E = Symbol.for;
    q = E("react.portal");
    r = E("react.fragment");
    u = E("react.strict_mode");
    z = E("react.profiler");
    B = E("react.provider");
    ba = E("react.context");
    ca = E("react.forward_ref");
    D = E("react.suspense");
    da = E("react.suspense_list");
    ea = E("react.memo");
    fa = E("react.lazy");
    ha = E("react.block");
    ia = E("react.fundamental");
    ja = E("react.scope");
    ka = E("react.debug_trace_mode");
    la = E("react.legacy_hidden");
  }
  var E;
  function F(a) {
    if (a == null)
      return null;
    if (typeof a === "function")
      return a.displayName || a.name || null;
    if (typeof a === "string")
      return a;
    switch (a) {
      case r:
        return "Fragment";
      case q:
        return "Portal";
      case z:
        return "Profiler";
      case u:
        return "StrictMode";
      case D:
        return "Suspense";
      case da:
        return "SuspenseList";
    }
    if (typeof a === "object")
      switch (a.$$typeof) {
        case ba:
          return (a.displayName || "Context") + ".Consumer";
        case B:
          return (a._context.displayName || "Context") + ".Provider";
        case ca:
          var b = a.render;
          b = b.displayName || b.name || "";
          return a.displayName || (b !== "" ? "ForwardRef(" + b + ")" : "ForwardRef");
        case ea:
          return F(a.type);
        case ha:
          return F(a._render);
        case fa:
          b = a._payload;
          a = a._init;
          try {
            return F(a(b));
          } catch (c) {
          }
      }
    return null;
  }
  var ma = n.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED;
  var na = {};
  function I(a, b) {
    for (var c = a._threadCount | 0; c <= b; c++)
      a[c] = a._currentValue2, a._threadCount = c + 1;
  }
  function oa(a, b, c, d) {
    if (d && (d = a.contextType, typeof d === "object" && d !== null))
      return I(d, c), d[c];
    if (a = a.contextTypes) {
      c = {};
      for (var f in a)
        c[f] = b[f];
      b = c;
    } else
      b = na;
    return b;
  }
  for (var J = new Uint16Array(16), K = 0; 15 > K; K++)
    J[K] = K + 1;
  J[15] = 0;
  var pa = /^[:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD][:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD\-.0-9\u00B7\u0300-\u036F\u203F-\u2040]*$/;
  var qa = Object.prototype.hasOwnProperty;
  var ra = {};
  var sa = {};
  function ta(a) {
    if (qa.call(sa, a))
      return true;
    if (qa.call(ra, a))
      return false;
    if (pa.test(a))
      return sa[a] = true;
    ra[a] = true;
    return false;
  }
  function ua(a, b, c, d) {
    if (c !== null && c.type === 0)
      return false;
    switch (typeof b) {
      case "function":
      case "symbol":
        return true;
      case "boolean":
        if (d)
          return false;
        if (c !== null)
          return !c.acceptsBooleans;
        a = a.toLowerCase().slice(0, 5);
        return a !== "data-" && a !== "aria-";
      default:
        return false;
    }
  }
  function va(a, b, c, d) {
    if (b === null || typeof b === "undefined" || ua(a, b, c, d))
      return true;
    if (d)
      return false;
    if (c !== null)
      switch (c.type) {
        case 3:
          return !b;
        case 4:
          return b === false;
        case 5:
          return isNaN(b);
        case 6:
          return isNaN(b) || 1 > b;
      }
    return false;
  }
  function M(a, b, c, d, f, h, t) {
    this.acceptsBooleans = b === 2 || b === 3 || b === 4;
    this.attributeName = d;
    this.attributeNamespace = f;
    this.mustUseProperty = c;
    this.propertyName = a;
    this.type = b;
    this.sanitizeURL = h;
    this.removeEmptyString = t;
  }
  var N = {};
  "children dangerouslySetInnerHTML defaultValue defaultChecked innerHTML suppressContentEditableWarning suppressHydrationWarning style".split(" ").forEach(function(a) {
    N[a] = new M(a, 0, false, a, null, false, false);
  });
  [["acceptCharset", "accept-charset"], ["className", "class"], ["htmlFor", "for"], ["httpEquiv", "http-equiv"]].forEach(function(a) {
    var b = a[0];
    N[b] = new M(b, 1, false, a[1], null, false, false);
  });
  ["contentEditable", "draggable", "spellCheck", "value"].forEach(function(a) {
    N[a] = new M(a, 2, false, a.toLowerCase(), null, false, false);
  });
  ["autoReverse", "externalResourcesRequired", "focusable", "preserveAlpha"].forEach(function(a) {
    N[a] = new M(a, 2, false, a, null, false, false);
  });
  "allowFullScreen async autoFocus autoPlay controls default defer disabled disablePictureInPicture disableRemotePlayback formNoValidate hidden loop noModule noValidate open playsInline readOnly required reversed scoped seamless itemScope".split(" ").forEach(function(a) {
    N[a] = new M(a, 3, false, a.toLowerCase(), null, false, false);
  });
  ["checked", "multiple", "muted", "selected"].forEach(function(a) {
    N[a] = new M(a, 3, true, a, null, false, false);
  });
  ["capture", "download"].forEach(function(a) {
    N[a] = new M(a, 4, false, a, null, false, false);
  });
  ["cols", "rows", "size", "span"].forEach(function(a) {
    N[a] = new M(a, 6, false, a, null, false, false);
  });
  ["rowSpan", "start"].forEach(function(a) {
    N[a] = new M(a, 5, false, a.toLowerCase(), null, false, false);
  });
  var wa = /[\-:]([a-z])/g;
  function xa(a) {
    return a[1].toUpperCase();
  }
  "accent-height alignment-baseline arabic-form baseline-shift cap-height clip-path clip-rule color-interpolation color-interpolation-filters color-profile color-rendering dominant-baseline enable-background fill-opacity fill-rule flood-color flood-opacity font-family font-size font-size-adjust font-stretch font-style font-variant font-weight glyph-name glyph-orientation-horizontal glyph-orientation-vertical horiz-adv-x horiz-origin-x image-rendering letter-spacing lighting-color marker-end marker-mid marker-start overline-position overline-thickness paint-order panose-1 pointer-events rendering-intent shape-rendering stop-color stop-opacity strikethrough-position strikethrough-thickness stroke-dasharray stroke-dashoffset stroke-linecap stroke-linejoin stroke-miterlimit stroke-opacity stroke-width text-anchor text-decoration text-rendering underline-position underline-thickness unicode-bidi unicode-range units-per-em v-alphabetic v-hanging v-ideographic v-mathematical vector-effect vert-adv-y vert-origin-x vert-origin-y word-spacing writing-mode xmlns:xlink x-height".split(" ").forEach(function(a) {
    var b = a.replace(wa, xa);
    N[b] = new M(b, 1, false, a, null, false, false);
  });
  "xlink:actuate xlink:arcrole xlink:role xlink:show xlink:title xlink:type".split(" ").forEach(function(a) {
    var b = a.replace(wa, xa);
    N[b] = new M(b, 1, false, a, "http://www.w3.org/1999/xlink", false, false);
  });
  ["xml:base", "xml:lang", "xml:space"].forEach(function(a) {
    var b = a.replace(wa, xa);
    N[b] = new M(b, 1, false, a, "http://www.w3.org/XML/1998/namespace", false, false);
  });
  ["tabIndex", "crossOrigin"].forEach(function(a) {
    N[a] = new M(a, 1, false, a.toLowerCase(), null, false, false);
  });
  N.xlinkHref = new M("xlinkHref", 1, false, "xlink:href", "http://www.w3.org/1999/xlink", true, false);
  ["src", "href", "action", "formAction"].forEach(function(a) {
    N[a] = new M(a, 1, false, a.toLowerCase(), null, true, true);
  });
  var ya = /["'&<>]/;
  function O(a) {
    if (typeof a === "boolean" || typeof a === "number")
      return "" + a;
    a = "" + a;
    var b = ya.exec(a);
    if (b) {
      var c = "", d, f = 0;
      for (d = b.index; d < a.length; d++) {
        switch (a.charCodeAt(d)) {
          case 34:
            b = "&quot;";
            break;
          case 38:
            b = "&amp;";
            break;
          case 39:
            b = "&#x27;";
            break;
          case 60:
            b = "&lt;";
            break;
          case 62:
            b = "&gt;";
            break;
          default:
            continue;
        }
        f !== d && (c += a.substring(f, d));
        f = d + 1;
        c += b;
      }
      a = f !== d ? c + a.substring(f, d) : c;
    }
    return a;
  }
  function za(a, b) {
    var c = N.hasOwnProperty(a) ? N[a] : null;
    var d;
    if (d = a !== "style")
      d = c !== null ? c.type === 0 : !(2 < a.length) || a[0] !== "o" && a[0] !== "O" || a[1] !== "n" && a[1] !== "N" ? false : true;
    if (d || va(a, b, c, false))
      return "";
    if (c !== null) {
      a = c.attributeName;
      d = c.type;
      if (d === 3 || d === 4 && b === true)
        return a + '=""';
      c.sanitizeURL && (b = "" + b);
      return a + '="' + (O(b) + '"');
    }
    return ta(a) ? a + '="' + (O(b) + '"') : "";
  }
  function Aa(a, b) {
    return a === b && (a !== 0 || 1 / a === 1 / b) || a !== a && b !== b;
  }
  var Ba = typeof Object.is === "function" ? Object.is : Aa;
  var P = null;
  var Q = null;
  var R = null;
  var S = false;
  var T = false;
  var U = null;
  var V = 0;
  function W() {
    if (P === null)
      throw Error(p(321));
    return P;
  }
  function Ca() {
    if (0 < V)
      throw Error(p(312));
    return {memoizedState: null, queue: null, next: null};
  }
  function Da() {
    R === null ? Q === null ? (S = false, Q = R = Ca()) : (S = true, R = Q) : R.next === null ? (S = false, R = R.next = Ca()) : (S = true, R = R.next);
    return R;
  }
  function Ea(a, b, c, d) {
    for (; T; )
      T = false, V += 1, R = null, c = a(b, d);
    Fa();
    return c;
  }
  function Fa() {
    P = null;
    T = false;
    Q = null;
    V = 0;
    R = U = null;
  }
  function Ga(a, b) {
    return typeof b === "function" ? b(a) : b;
  }
  function Ha(a, b, c) {
    P = W();
    R = Da();
    if (S) {
      var d = R.queue;
      b = d.dispatch;
      if (U !== null && (c = U.get(d), c !== void 0)) {
        U.delete(d);
        d = R.memoizedState;
        do
          d = a(d, c.action), c = c.next;
        while (c !== null);
        R.memoizedState = d;
        return [d, b];
      }
      return [R.memoizedState, b];
    }
    a = a === Ga ? typeof b === "function" ? b() : b : c !== void 0 ? c(b) : b;
    R.memoizedState = a;
    a = R.queue = {last: null, dispatch: null};
    a = a.dispatch = Ia.bind(null, P, a);
    return [R.memoizedState, a];
  }
  function Ja(a, b) {
    P = W();
    R = Da();
    b = b === void 0 ? null : b;
    if (R !== null) {
      var c = R.memoizedState;
      if (c !== null && b !== null) {
        var d = c[1];
        a:
          if (d === null)
            d = false;
          else {
            for (var f = 0; f < d.length && f < b.length; f++)
              if (!Ba(b[f], d[f])) {
                d = false;
                break a;
              }
            d = true;
          }
        if (d)
          return c[0];
      }
    }
    a = a();
    R.memoizedState = [a, b];
    return a;
  }
  function Ia(a, b, c) {
    if (!(25 > V))
      throw Error(p(301));
    if (a === P)
      if (T = true, a = {action: c, next: null}, U === null && (U = new Map()), c = U.get(b), c === void 0)
        U.set(b, a);
      else {
        for (b = c; b.next !== null; )
          b = b.next;
        b.next = a;
      }
  }
  function Ka() {
  }
  var X = null;
  var La = {readContext: function(a) {
    var b = X.threadID;
    I(a, b);
    return a[b];
  }, useContext: function(a) {
    W();
    var b = X.threadID;
    I(a, b);
    return a[b];
  }, useMemo: Ja, useReducer: Ha, useRef: function(a) {
    P = W();
    R = Da();
    var b = R.memoizedState;
    return b === null ? (a = {current: a}, R.memoizedState = a) : b;
  }, useState: function(a) {
    return Ha(Ga, a);
  }, useLayoutEffect: function() {
  }, useCallback: function(a, b) {
    return Ja(function() {
      return a;
    }, b);
  }, useImperativeHandle: Ka, useEffect: Ka, useDebugValue: Ka, useDeferredValue: function(a) {
    W();
    return a;
  }, useTransition: function() {
    W();
    return [function(a) {
      a();
    }, false];
  }, useOpaqueIdentifier: function() {
    return (X.identifierPrefix || "") + "R:" + (X.uniqueID++).toString(36);
  }, useMutableSource: function(a, b) {
    W();
    return b(a._source);
  }};
  var Ma = {html: "http://www.w3.org/1999/xhtml", mathml: "http://www.w3.org/1998/Math/MathML", svg: "http://www.w3.org/2000/svg"};
  function Na(a) {
    switch (a) {
      case "svg":
        return "http://www.w3.org/2000/svg";
      case "math":
        return "http://www.w3.org/1998/Math/MathML";
      default:
        return "http://www.w3.org/1999/xhtml";
    }
  }
  var Oa = {area: true, base: true, br: true, col: true, embed: true, hr: true, img: true, input: true, keygen: true, link: true, meta: true, param: true, source: true, track: true, wbr: true};
  var Pa = l({menuitem: true}, Oa);
  var Y = {
    animationIterationCount: true,
    borderImageOutset: true,
    borderImageSlice: true,
    borderImageWidth: true,
    boxFlex: true,
    boxFlexGroup: true,
    boxOrdinalGroup: true,
    columnCount: true,
    columns: true,
    flex: true,
    flexGrow: true,
    flexPositive: true,
    flexShrink: true,
    flexNegative: true,
    flexOrder: true,
    gridArea: true,
    gridRow: true,
    gridRowEnd: true,
    gridRowSpan: true,
    gridRowStart: true,
    gridColumn: true,
    gridColumnEnd: true,
    gridColumnSpan: true,
    gridColumnStart: true,
    fontWeight: true,
    lineClamp: true,
    lineHeight: true,
    opacity: true,
    order: true,
    orphans: true,
    tabSize: true,
    widows: true,
    zIndex: true,
    zoom: true,
    fillOpacity: true,
    floodOpacity: true,
    stopOpacity: true,
    strokeDasharray: true,
    strokeDashoffset: true,
    strokeMiterlimit: true,
    strokeOpacity: true,
    strokeWidth: true
  };
  var Qa = ["Webkit", "ms", "Moz", "O"];
  Object.keys(Y).forEach(function(a) {
    Qa.forEach(function(b) {
      b = b + a.charAt(0).toUpperCase() + a.substring(1);
      Y[b] = Y[a];
    });
  });
  var Ra = /([A-Z])/g;
  var Sa = /^ms-/;
  var Z = n.Children.toArray;
  var Ta = ma.ReactCurrentDispatcher;
  var Ua = {listing: true, pre: true, textarea: true};
  var Va = /^[a-zA-Z][a-zA-Z:_\.\-\d]*$/;
  var Wa = {};
  var Xa = {};
  function Ya(a) {
    if (a === void 0 || a === null)
      return a;
    var b = "";
    n.Children.forEach(a, function(a2) {
      a2 != null && (b += a2);
    });
    return b;
  }
  var Za = Object.prototype.hasOwnProperty;
  var $a = {children: null, dangerouslySetInnerHTML: null, suppressContentEditableWarning: null, suppressHydrationWarning: null};
  function ab(a, b) {
    if (a === void 0)
      throw Error(p(152, F(b) || "Component"));
  }
  function bb(a, b, c) {
    function d(d2, h2) {
      var e = h2.prototype && h2.prototype.isReactComponent, f2 = oa(h2, b, c, e), t = [], g = false, m = {isMounted: function() {
        return false;
      }, enqueueForceUpdate: function() {
        if (t === null)
          return null;
      }, enqueueReplaceState: function(a2, b2) {
        g = true;
        t = [b2];
      }, enqueueSetState: function(a2, b2) {
        if (t === null)
          return null;
        t.push(b2);
      }};
      if (e) {
        if (e = new h2(d2.props, f2, m), typeof h2.getDerivedStateFromProps === "function") {
          var k = h2.getDerivedStateFromProps.call(null, d2.props, e.state);
          k != null && (e.state = l({}, e.state, k));
        }
      } else if (P = {}, e = h2(d2.props, f2, m), e = Ea(h2, d2.props, e, f2), e == null || e.render == null) {
        a = e;
        ab(a, h2);
        return;
      }
      e.props = d2.props;
      e.context = f2;
      e.updater = m;
      m = e.state;
      m === void 0 && (e.state = m = null);
      if (typeof e.UNSAFE_componentWillMount === "function" || typeof e.componentWillMount === "function")
        if (typeof e.componentWillMount === "function" && typeof h2.getDerivedStateFromProps !== "function" && e.componentWillMount(), typeof e.UNSAFE_componentWillMount === "function" && typeof h2.getDerivedStateFromProps !== "function" && e.UNSAFE_componentWillMount(), t.length) {
          m = t;
          var v = g;
          t = null;
          g = false;
          if (v && m.length === 1)
            e.state = m[0];
          else {
            k = v ? m[0] : e.state;
            var H = true;
            for (v = v ? 1 : 0; v < m.length; v++) {
              var x = m[v];
              x = typeof x === "function" ? x.call(e, k, d2.props, f2) : x;
              x != null && (H ? (H = false, k = l({}, k, x)) : l(k, x));
            }
            e.state = k;
          }
        } else
          t = null;
      a = e.render();
      ab(a, h2);
      if (typeof e.getChildContext === "function" && (d2 = h2.childContextTypes, typeof d2 === "object")) {
        var y = e.getChildContext();
        for (var A in y)
          if (!(A in d2))
            throw Error(p(108, F(h2) || "Unknown", A));
      }
      y && (b = l({}, b, y));
    }
    for (; n.isValidElement(a); ) {
      var f = a, h = f.type;
      if (typeof h !== "function")
        break;
      d(f, h);
    }
    return {child: a, context: b};
  }
  var cb = function() {
    function a(a2, b2, f) {
      n.isValidElement(a2) ? a2.type !== r ? a2 = [a2] : (a2 = a2.props.children, a2 = n.isValidElement(a2) ? [a2] : Z(a2)) : a2 = Z(a2);
      a2 = {type: null, domNamespace: Ma.html, children: a2, childIndex: 0, context: na, footer: ""};
      var c = J[0];
      if (c === 0) {
        var d = J;
        c = d.length;
        var g = 2 * c;
        if (!(65536 >= g))
          throw Error(p(304));
        var e = new Uint16Array(g);
        e.set(d);
        J = e;
        J[0] = c + 1;
        for (d = c; d < g - 1; d++)
          J[d] = d + 1;
        J[g - 1] = 0;
      } else
        J[0] = J[c];
      this.threadID = c;
      this.stack = [a2];
      this.exhausted = false;
      this.currentSelectValue = null;
      this.previousWasTextNode = false;
      this.makeStaticMarkup = b2;
      this.suspenseDepth = 0;
      this.contextIndex = -1;
      this.contextStack = [];
      this.contextValueStack = [];
      this.uniqueID = 0;
      this.identifierPrefix = f && f.identifierPrefix || "";
    }
    var b = a.prototype;
    b.destroy = function() {
      if (!this.exhausted) {
        this.exhausted = true;
        this.clearProviders();
        var a2 = this.threadID;
        J[a2] = J[0];
        J[0] = a2;
      }
    };
    b.pushProvider = function(a2) {
      var b2 = ++this.contextIndex, c = a2.type._context, h = this.threadID;
      I(c, h);
      var t = c[h];
      this.contextStack[b2] = c;
      this.contextValueStack[b2] = t;
      c[h] = a2.props.value;
    };
    b.popProvider = function() {
      var a2 = this.contextIndex, b2 = this.contextStack[a2], f = this.contextValueStack[a2];
      this.contextStack[a2] = null;
      this.contextValueStack[a2] = null;
      this.contextIndex--;
      b2[this.threadID] = f;
    };
    b.clearProviders = function() {
      for (var a2 = this.contextIndex; 0 <= a2; a2--)
        this.contextStack[a2][this.threadID] = this.contextValueStack[a2];
    };
    b.read = function(a2) {
      if (this.exhausted)
        return null;
      var b2 = X;
      X = this;
      var c = Ta.current;
      Ta.current = La;
      try {
        for (var h = [""], t = false; h[0].length < a2; ) {
          if (this.stack.length === 0) {
            this.exhausted = true;
            var g = this.threadID;
            J[g] = J[0];
            J[0] = g;
            break;
          }
          var e = this.stack[this.stack.length - 1];
          if (t || e.childIndex >= e.children.length) {
            var L = e.footer;
            L !== "" && (this.previousWasTextNode = false);
            this.stack.pop();
            if (e.type === "select")
              this.currentSelectValue = null;
            else if (e.type != null && e.type.type != null && e.type.type.$$typeof === B)
              this.popProvider(e.type);
            else if (e.type === D) {
              this.suspenseDepth--;
              var G = h.pop();
              if (t) {
                t = false;
                var C = e.fallbackFrame;
                if (!C)
                  throw Error(p(303));
                this.stack.push(C);
                h[this.suspenseDepth] += "<!--$!-->";
                continue;
              } else
                h[this.suspenseDepth] += G;
            }
            h[this.suspenseDepth] += L;
          } else {
            var m = e.children[e.childIndex++], k = "";
            try {
              k += this.render(m, e.context, e.domNamespace);
            } catch (v) {
              if (v != null && typeof v.then === "function")
                throw Error(p(294));
              throw v;
            } finally {
            }
            h.length <= this.suspenseDepth && h.push("");
            h[this.suspenseDepth] += k;
          }
        }
        return h[0];
      } finally {
        Ta.current = c, X = b2, Fa();
      }
    };
    b.render = function(a2, b2, f) {
      if (typeof a2 === "string" || typeof a2 === "number") {
        f = "" + a2;
        if (f === "")
          return "";
        if (this.makeStaticMarkup)
          return O(f);
        if (this.previousWasTextNode)
          return "<!-- -->" + O(f);
        this.previousWasTextNode = true;
        return O(f);
      }
      b2 = bb(a2, b2, this.threadID);
      a2 = b2.child;
      b2 = b2.context;
      if (a2 === null || a2 === false)
        return "";
      if (!n.isValidElement(a2)) {
        if (a2 != null && a2.$$typeof != null) {
          f = a2.$$typeof;
          if (f === q)
            throw Error(p(257));
          throw Error(p(258, f.toString()));
        }
        a2 = Z(a2);
        this.stack.push({type: null, domNamespace: f, children: a2, childIndex: 0, context: b2, footer: ""});
        return "";
      }
      var c = a2.type;
      if (typeof c === "string")
        return this.renderDOM(a2, b2, f);
      switch (c) {
        case la:
        case ka:
        case u:
        case z:
        case da:
        case r:
          return a2 = Z(a2.props.children), this.stack.push({
            type: null,
            domNamespace: f,
            children: a2,
            childIndex: 0,
            context: b2,
            footer: ""
          }), "";
        case D:
          throw Error(p(294));
        case ja:
          throw Error(p(343));
      }
      if (typeof c === "object" && c !== null)
        switch (c.$$typeof) {
          case ca:
            P = {};
            var d = c.render(a2.props, a2.ref);
            d = Ea(c.render, a2.props, d, a2.ref);
            d = Z(d);
            this.stack.push({type: null, domNamespace: f, children: d, childIndex: 0, context: b2, footer: ""});
            return "";
          case ea:
            return a2 = [n.createElement(c.type, l({ref: a2.ref}, a2.props))], this.stack.push({type: null, domNamespace: f, children: a2, childIndex: 0, context: b2, footer: ""}), "";
          case B:
            return c = Z(a2.props.children), f = {type: a2, domNamespace: f, children: c, childIndex: 0, context: b2, footer: ""}, this.pushProvider(a2), this.stack.push(f), "";
          case ba:
            c = a2.type;
            d = a2.props;
            var g = this.threadID;
            I(c, g);
            c = Z(d.children(c[g]));
            this.stack.push({type: a2, domNamespace: f, children: c, childIndex: 0, context: b2, footer: ""});
            return "";
          case ia:
            throw Error(p(338));
          case fa:
            return c = a2.type, d = c._init, c = d(c._payload), a2 = [n.createElement(c, l({ref: a2.ref}, a2.props))], this.stack.push({
              type: null,
              domNamespace: f,
              children: a2,
              childIndex: 0,
              context: b2,
              footer: ""
            }), "";
        }
      throw Error(p(130, c == null ? c : typeof c, ""));
    };
    b.renderDOM = function(a2, b2, f) {
      var c = a2.type.toLowerCase();
      f === Ma.html && Na(c);
      if (!Wa.hasOwnProperty(c)) {
        if (!Va.test(c))
          throw Error(p(65, c));
        Wa[c] = true;
      }
      var d = a2.props;
      if (c === "input")
        d = l({type: void 0}, d, {defaultChecked: void 0, defaultValue: void 0, value: d.value != null ? d.value : d.defaultValue, checked: d.checked != null ? d.checked : d.defaultChecked});
      else if (c === "textarea") {
        var g = d.value;
        if (g == null) {
          g = d.defaultValue;
          var e = d.children;
          if (e != null) {
            if (g != null)
              throw Error(p(92));
            if (Array.isArray(e)) {
              if (!(1 >= e.length))
                throw Error(p(93));
              e = e[0];
            }
            g = "" + e;
          }
          g == null && (g = "");
        }
        d = l({}, d, {value: void 0, children: "" + g});
      } else if (c === "select")
        this.currentSelectValue = d.value != null ? d.value : d.defaultValue, d = l({}, d, {value: void 0});
      else if (c === "option") {
        e = this.currentSelectValue;
        var L = Ya(d.children);
        if (e != null) {
          var G = d.value != null ? d.value + "" : L;
          g = false;
          if (Array.isArray(e))
            for (var C = 0; C < e.length; C++) {
              if ("" + e[C] === G) {
                g = true;
                break;
              }
            }
          else
            g = "" + e === G;
          d = l({selected: void 0, children: void 0}, d, {selected: g, children: L});
        }
      }
      if (g = d) {
        if (Pa[c] && (g.children != null || g.dangerouslySetInnerHTML != null))
          throw Error(p(137, c));
        if (g.dangerouslySetInnerHTML != null) {
          if (g.children != null)
            throw Error(p(60));
          if (!(typeof g.dangerouslySetInnerHTML === "object" && "__html" in g.dangerouslySetInnerHTML))
            throw Error(p(61));
        }
        if (g.style != null && typeof g.style !== "object")
          throw Error(p(62));
      }
      g = d;
      e = this.makeStaticMarkup;
      L = this.stack.length === 1;
      G = "<" + a2.type;
      b:
        if (c.indexOf("-") === -1)
          C = typeof g.is === "string";
        else
          switch (c) {
            case "annotation-xml":
            case "color-profile":
            case "font-face":
            case "font-face-src":
            case "font-face-uri":
            case "font-face-format":
            case "font-face-name":
            case "missing-glyph":
              C = false;
              break b;
            default:
              C = true;
          }
      for (w in g)
        if (Za.call(g, w)) {
          var m = g[w];
          if (m != null) {
            if (w === "style") {
              var k = void 0, v = "", H = "";
              for (k in m)
                if (m.hasOwnProperty(k)) {
                  var x = k.indexOf("--") === 0, y = m[k];
                  if (y != null) {
                    if (x)
                      var A = k;
                    else if (A = k, Xa.hasOwnProperty(A))
                      A = Xa[A];
                    else {
                      var eb = A.replace(Ra, "-$1").toLowerCase().replace(Sa, "-ms-");
                      A = Xa[A] = eb;
                    }
                    v += H + A + ":";
                    H = k;
                    x = y == null || typeof y === "boolean" || y === "" ? "" : x || typeof y !== "number" || y === 0 || Y.hasOwnProperty(H) && Y[H] ? ("" + y).trim() : y + "px";
                    v += x;
                    H = ";";
                  }
                }
              m = v || null;
            }
            k = null;
            C ? $a.hasOwnProperty(w) || (k = w, k = ta(k) && m != null ? k + '="' + (O(m) + '"') : "") : k = za(w, m);
            k && (G += " " + k);
          }
        }
      e || L && (G += ' data-reactroot=""');
      var w = G;
      g = "";
      Oa.hasOwnProperty(c) ? w += "/>" : (w += ">", g = "</" + a2.type + ">");
      a: {
        e = d.dangerouslySetInnerHTML;
        if (e != null) {
          if (e.__html != null) {
            e = e.__html;
            break a;
          }
        } else if (e = d.children, typeof e === "string" || typeof e === "number") {
          e = O(e);
          break a;
        }
        e = null;
      }
      e != null ? (d = [], Ua.hasOwnProperty(c) && e.charAt(0) === "\n" && (w += "\n"), w += e) : d = Z(d.children);
      a2 = a2.type;
      f = f == null || f === "http://www.w3.org/1999/xhtml" ? Na(a2) : f === "http://www.w3.org/2000/svg" && a2 === "foreignObject" ? "http://www.w3.org/1999/xhtml" : f;
      this.stack.push({domNamespace: f, type: c, children: d, childIndex: 0, context: b2, footer: g});
      this.previousWasTextNode = false;
      return w;
    };
    return a;
  }();
  function db(a, b) {
    a.prototype = Object.create(b.prototype);
    a.prototype.constructor = a;
    a.__proto__ = b;
  }
  var fb = function(a) {
    function b(b2, c2, h) {
      var d = a.call(this, {}) || this;
      d.partialRenderer = new cb(b2, c2, h);
      return d;
    }
    db(b, a);
    var c = b.prototype;
    c._destroy = function(a2, b2) {
      this.partialRenderer.destroy();
      b2(a2);
    };
    c._read = function(a2) {
      try {
        this.push(this.partialRenderer.read(a2));
      } catch (f) {
        this.destroy(f);
      }
    };
    return b;
  }(aa.Readable);
  exports2.renderToNodeStream = function(a, b) {
    return new fb(a, false, b);
  };
  exports2.renderToStaticMarkup = function(a, b) {
    a = new cb(a, true, b);
    try {
      return a.read(Infinity);
    } finally {
      a.destroy();
    }
  };
  exports2.renderToStaticNodeStream = function(a, b) {
    return new fb(a, true, b);
  };
  exports2.renderToString = function(a, b) {
    a = new cb(a, false, b);
    try {
      return a.read(Infinity);
    } finally {
      a.destroy();
    }
  };
  exports2.version = "17.0.1";
});

// node_modules/react-dom/server.node.js
var require_server_node = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_react_dom_server_node_production_min();
  } else {
    module2.exports = null;
  }
});

// node_modules/react-dom/server.js
var require_server = __commonJS((exports2, module2) => {
  "use strict";
  module2.exports = require_server_node();
});

// node_modules/marked/src/defaults.js
var require_defaults = __commonJS((exports2, module2) => {
  function getDefaults() {
    return {
      baseUrl: null,
      breaks: false,
      gfm: true,
      headerIds: true,
      headerPrefix: "",
      highlight: null,
      langPrefix: "language-",
      mangle: true,
      pedantic: false,
      renderer: null,
      sanitize: false,
      sanitizer: null,
      silent: false,
      smartLists: false,
      smartypants: false,
      tokenizer: null,
      walkTokens: null,
      xhtml: false
    };
  }
  function changeDefaults(newDefaults) {
    module2.exports.defaults = newDefaults;
  }
  module2.exports = {
    defaults: getDefaults(),
    getDefaults,
    changeDefaults
  };
});

// node_modules/marked/src/helpers.js
var require_helpers = __commonJS((exports2, module2) => {
  var escapeTest = /[&<>"']/;
  var escapeReplace = /[&<>"']/g;
  var escapeTestNoEncode = /[<>"']|&(?!#?\w+;)/;
  var escapeReplaceNoEncode = /[<>"']|&(?!#?\w+;)/g;
  var escapeReplacements = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#39;"
  };
  var getEscapeReplacement = (ch) => escapeReplacements[ch];
  function escape(html, encode) {
    if (encode) {
      if (escapeTest.test(html)) {
        return html.replace(escapeReplace, getEscapeReplacement);
      }
    } else {
      if (escapeTestNoEncode.test(html)) {
        return html.replace(escapeReplaceNoEncode, getEscapeReplacement);
      }
    }
    return html;
  }
  var unescapeTest = /&(#(?:\d+)|(?:#x[0-9A-Fa-f]+)|(?:\w+));?/ig;
  function unescape(html) {
    return html.replace(unescapeTest, (_, n) => {
      n = n.toLowerCase();
      if (n === "colon")
        return ":";
      if (n.charAt(0) === "#") {
        return n.charAt(1) === "x" ? String.fromCharCode(parseInt(n.substring(2), 16)) : String.fromCharCode(+n.substring(1));
      }
      return "";
    });
  }
  var caret = /(^|[^\[])\^/g;
  function edit(regex, opt) {
    regex = regex.source || regex;
    opt = opt || "";
    const obj = {
      replace: (name, val) => {
        val = val.source || val;
        val = val.replace(caret, "$1");
        regex = regex.replace(name, val);
        return obj;
      },
      getRegex: () => {
        return new RegExp(regex, opt);
      }
    };
    return obj;
  }
  var nonWordAndColonTest = /[^\w:]/g;
  var originIndependentUrl = /^$|^[a-z][a-z0-9+.-]*:|^[?#]/i;
  function cleanUrl(sanitize, base, href) {
    if (sanitize) {
      let prot;
      try {
        prot = decodeURIComponent(unescape(href)).replace(nonWordAndColonTest, "").toLowerCase();
      } catch (e) {
        return null;
      }
      if (prot.indexOf("javascript:") === 0 || prot.indexOf("vbscript:") === 0 || prot.indexOf("data:") === 0) {
        return null;
      }
    }
    if (base && !originIndependentUrl.test(href)) {
      href = resolveUrl(base, href);
    }
    try {
      href = encodeURI(href).replace(/%25/g, "%");
    } catch (e) {
      return null;
    }
    return href;
  }
  var baseUrls = {};
  var justDomain = /^[^:]+:\/*[^/]*$/;
  var protocol = /^([^:]+:)[\s\S]*$/;
  var domain = /^([^:]+:\/*[^/]*)[\s\S]*$/;
  function resolveUrl(base, href) {
    if (!baseUrls[" " + base]) {
      if (justDomain.test(base)) {
        baseUrls[" " + base] = base + "/";
      } else {
        baseUrls[" " + base] = rtrim(base, "/", true);
      }
    }
    base = baseUrls[" " + base];
    const relativeBase = base.indexOf(":") === -1;
    if (href.substring(0, 2) === "//") {
      if (relativeBase) {
        return href;
      }
      return base.replace(protocol, "$1") + href;
    } else if (href.charAt(0) === "/") {
      if (relativeBase) {
        return href;
      }
      return base.replace(domain, "$1") + href;
    } else {
      return base + href;
    }
  }
  var noopTest = {exec: function noopTest2() {
  }};
  function merge(obj) {
    let i = 1, target, key;
    for (; i < arguments.length; i++) {
      target = arguments[i];
      for (key in target) {
        if (Object.prototype.hasOwnProperty.call(target, key)) {
          obj[key] = target[key];
        }
      }
    }
    return obj;
  }
  function splitCells(tableRow, count) {
    const row = tableRow.replace(/\|/g, (match, offset, str) => {
      let escaped = false, curr = offset;
      while (--curr >= 0 && str[curr] === "\\")
        escaped = !escaped;
      if (escaped) {
        return "|";
      } else {
        return " |";
      }
    }), cells = row.split(/ \|/);
    let i = 0;
    if (cells.length > count) {
      cells.splice(count);
    } else {
      while (cells.length < count)
        cells.push("");
    }
    for (; i < cells.length; i++) {
      cells[i] = cells[i].trim().replace(/\\\|/g, "|");
    }
    return cells;
  }
  function rtrim(str, c, invert) {
    const l = str.length;
    if (l === 0) {
      return "";
    }
    let suffLen = 0;
    while (suffLen < l) {
      const currChar = str.charAt(l - suffLen - 1);
      if (currChar === c && !invert) {
        suffLen++;
      } else if (currChar !== c && invert) {
        suffLen++;
      } else {
        break;
      }
    }
    return str.substr(0, l - suffLen);
  }
  function findClosingBracket(str, b) {
    if (str.indexOf(b[1]) === -1) {
      return -1;
    }
    const l = str.length;
    let level = 0, i = 0;
    for (; i < l; i++) {
      if (str[i] === "\\") {
        i++;
      } else if (str[i] === b[0]) {
        level++;
      } else if (str[i] === b[1]) {
        level--;
        if (level < 0) {
          return i;
        }
      }
    }
    return -1;
  }
  function checkSanitizeDeprecation(opt) {
    if (opt && opt.sanitize && !opt.silent) {
      console.warn("marked(): sanitize and sanitizer parameters are deprecated since version 0.7.0, should not be used and will be removed in the future. Read more here: https://marked.js.org/#/USING_ADVANCED.md#options");
    }
  }
  function repeatString(pattern, count) {
    if (count < 1) {
      return "";
    }
    let result = "";
    while (count > 1) {
      if (count & 1) {
        result += pattern;
      }
      count >>= 1;
      pattern += pattern;
    }
    return result + pattern;
  }
  module2.exports = {
    escape,
    unescape,
    edit,
    cleanUrl,
    resolveUrl,
    noopTest,
    merge,
    splitCells,
    rtrim,
    findClosingBracket,
    checkSanitizeDeprecation,
    repeatString
  };
});

// node_modules/marked/src/Tokenizer.js
var require_Tokenizer = __commonJS((exports2, module2) => {
  var {defaults} = require_defaults();
  var {
    rtrim,
    splitCells,
    escape,
    findClosingBracket
  } = require_helpers();
  function outputLink(cap, link2, raw) {
    const href = link2.href;
    const title = link2.title ? escape(link2.title) : null;
    const text = cap[1].replace(/\\([\[\]])/g, "$1");
    if (cap[0].charAt(0) !== "!") {
      return {
        type: "link",
        raw,
        href,
        title,
        text
      };
    } else {
      return {
        type: "image",
        raw,
        href,
        title,
        text: escape(text)
      };
    }
  }
  function indentCodeCompensation(raw, text) {
    const matchIndentToCode = raw.match(/^(\s+)(?:```)/);
    if (matchIndentToCode === null) {
      return text;
    }
    const indentToCode = matchIndentToCode[1];
    return text.split("\n").map((node) => {
      const matchIndentInNode = node.match(/^\s+/);
      if (matchIndentInNode === null) {
        return node;
      }
      const [indentInNode] = matchIndentInNode;
      if (indentInNode.length >= indentToCode.length) {
        return node.slice(indentToCode.length);
      }
      return node;
    }).join("\n");
  }
  module2.exports = class Tokenizer {
    constructor(options) {
      this.options = options || defaults;
    }
    space(src) {
      const cap = this.rules.block.newline.exec(src);
      if (cap) {
        if (cap[0].length > 1) {
          return {
            type: "space",
            raw: cap[0]
          };
        }
        return {raw: "\n"};
      }
    }
    code(src) {
      const cap = this.rules.block.code.exec(src);
      if (cap) {
        const text = cap[0].replace(/^ {1,4}/gm, "");
        return {
          type: "code",
          raw: cap[0],
          codeBlockStyle: "indented",
          text: !this.options.pedantic ? rtrim(text, "\n") : text
        };
      }
    }
    fences(src) {
      const cap = this.rules.block.fences.exec(src);
      if (cap) {
        const raw = cap[0];
        const text = indentCodeCompensation(raw, cap[3] || "");
        return {
          type: "code",
          raw,
          lang: cap[2] ? cap[2].trim() : cap[2],
          text
        };
      }
    }
    heading(src) {
      const cap = this.rules.block.heading.exec(src);
      if (cap) {
        let text = cap[2].trim();
        if (/#$/.test(text)) {
          const trimmed = rtrim(text, "#");
          if (this.options.pedantic) {
            text = trimmed.trim();
          } else if (!trimmed || / $/.test(trimmed)) {
            text = trimmed.trim();
          }
        }
        return {
          type: "heading",
          raw: cap[0],
          depth: cap[1].length,
          text
        };
      }
    }
    nptable(src) {
      const cap = this.rules.block.nptable.exec(src);
      if (cap) {
        const item = {
          type: "table",
          header: splitCells(cap[1].replace(/^ *| *\| *$/g, "")),
          align: cap[2].replace(/^ *|\| *$/g, "").split(/ *\| */),
          cells: cap[3] ? cap[3].replace(/\n$/, "").split("\n") : [],
          raw: cap[0]
        };
        if (item.header.length === item.align.length) {
          let l = item.align.length;
          let i;
          for (i = 0; i < l; i++) {
            if (/^ *-+: *$/.test(item.align[i])) {
              item.align[i] = "right";
            } else if (/^ *:-+: *$/.test(item.align[i])) {
              item.align[i] = "center";
            } else if (/^ *:-+ *$/.test(item.align[i])) {
              item.align[i] = "left";
            } else {
              item.align[i] = null;
            }
          }
          l = item.cells.length;
          for (i = 0; i < l; i++) {
            item.cells[i] = splitCells(item.cells[i], item.header.length);
          }
          return item;
        }
      }
    }
    hr(src) {
      const cap = this.rules.block.hr.exec(src);
      if (cap) {
        return {
          type: "hr",
          raw: cap[0]
        };
      }
    }
    blockquote(src) {
      const cap = this.rules.block.blockquote.exec(src);
      if (cap) {
        const text = cap[0].replace(/^ *> ?/gm, "");
        return {
          type: "blockquote",
          raw: cap[0],
          text
        };
      }
    }
    list(src) {
      const cap = this.rules.block.list.exec(src);
      if (cap) {
        let raw = cap[0];
        const bull = cap[2];
        const isordered = bull.length > 1;
        const list = {
          type: "list",
          raw,
          ordered: isordered,
          start: isordered ? +bull.slice(0, -1) : "",
          loose: false,
          items: []
        };
        const itemMatch = cap[0].match(this.rules.block.item);
        let next = false, item, space, bcurr, bnext, addBack, loose, istask, ischecked, endMatch;
        let l = itemMatch.length;
        bcurr = this.rules.block.listItemStart.exec(itemMatch[0]);
        for (let i = 0; i < l; i++) {
          item = itemMatch[i];
          raw = item;
          if (!this.options.pedantic) {
            endMatch = item.match(new RegExp("\\n\\s*\\n {0," + (bcurr[0].length - 1) + "}\\S"));
            if (endMatch) {
              addBack = item.length - endMatch.index + itemMatch.slice(i + 1).join("\n").length;
              list.raw = list.raw.substring(0, list.raw.length - addBack);
              item = item.substring(0, endMatch.index);
              raw = item;
              l = i + 1;
            }
          }
          if (i !== l - 1) {
            bnext = this.rules.block.listItemStart.exec(itemMatch[i + 1]);
            if (!this.options.pedantic ? bnext[1].length >= bcurr[0].length || bnext[1].length > 3 : bnext[1].length > bcurr[1].length) {
              itemMatch.splice(i, 2, itemMatch[i] + (!this.options.pedantic && bnext[1].length < bcurr[0].length && !itemMatch[i].match(/\n$/) ? "" : "\n") + itemMatch[i + 1]);
              i--;
              l--;
              continue;
            } else if (!this.options.pedantic || this.options.smartLists ? bnext[2][bnext[2].length - 1] !== bull[bull.length - 1] : isordered === (bnext[2].length === 1)) {
              addBack = itemMatch.slice(i + 1).join("\n").length;
              list.raw = list.raw.substring(0, list.raw.length - addBack);
              i = l - 1;
            }
            bcurr = bnext;
          }
          space = item.length;
          item = item.replace(/^ *([*+-]|\d+[.)]) ?/, "");
          if (~item.indexOf("\n ")) {
            space -= item.length;
            item = !this.options.pedantic ? item.replace(new RegExp("^ {1," + space + "}", "gm"), "") : item.replace(/^ {1,4}/gm, "");
          }
          item = rtrim(item, "\n");
          if (i !== l - 1) {
            raw = raw + "\n";
          }
          loose = next || /\n\n(?!\s*$)/.test(raw);
          if (i !== l - 1) {
            next = raw.slice(-2) === "\n\n";
            if (!loose)
              loose = next;
          }
          if (loose) {
            list.loose = true;
          }
          if (this.options.gfm) {
            istask = /^\[[ xX]\] /.test(item);
            ischecked = void 0;
            if (istask) {
              ischecked = item[1] !== " ";
              item = item.replace(/^\[[ xX]\] +/, "");
            }
          }
          list.items.push({
            type: "list_item",
            raw,
            task: istask,
            checked: ischecked,
            loose,
            text: item
          });
        }
        return list;
      }
    }
    html(src) {
      const cap = this.rules.block.html.exec(src);
      if (cap) {
        return {
          type: this.options.sanitize ? "paragraph" : "html",
          raw: cap[0],
          pre: !this.options.sanitizer && (cap[1] === "pre" || cap[1] === "script" || cap[1] === "style"),
          text: this.options.sanitize ? this.options.sanitizer ? this.options.sanitizer(cap[0]) : escape(cap[0]) : cap[0]
        };
      }
    }
    def(src) {
      const cap = this.rules.block.def.exec(src);
      if (cap) {
        if (cap[3])
          cap[3] = cap[3].substring(1, cap[3].length - 1);
        const tag = cap[1].toLowerCase().replace(/\s+/g, " ");
        return {
          tag,
          raw: cap[0],
          href: cap[2],
          title: cap[3]
        };
      }
    }
    table(src) {
      const cap = this.rules.block.table.exec(src);
      if (cap) {
        const item = {
          type: "table",
          header: splitCells(cap[1].replace(/^ *| *\| *$/g, "")),
          align: cap[2].replace(/^ *|\| *$/g, "").split(/ *\| */),
          cells: cap[3] ? cap[3].replace(/\n$/, "").split("\n") : []
        };
        if (item.header.length === item.align.length) {
          item.raw = cap[0];
          let l = item.align.length;
          let i;
          for (i = 0; i < l; i++) {
            if (/^ *-+: *$/.test(item.align[i])) {
              item.align[i] = "right";
            } else if (/^ *:-+: *$/.test(item.align[i])) {
              item.align[i] = "center";
            } else if (/^ *:-+ *$/.test(item.align[i])) {
              item.align[i] = "left";
            } else {
              item.align[i] = null;
            }
          }
          l = item.cells.length;
          for (i = 0; i < l; i++) {
            item.cells[i] = splitCells(item.cells[i].replace(/^ *\| *| *\| *$/g, ""), item.header.length);
          }
          return item;
        }
      }
    }
    lheading(src) {
      const cap = this.rules.block.lheading.exec(src);
      if (cap) {
        return {
          type: "heading",
          raw: cap[0],
          depth: cap[2].charAt(0) === "=" ? 1 : 2,
          text: cap[1]
        };
      }
    }
    paragraph(src) {
      const cap = this.rules.block.paragraph.exec(src);
      if (cap) {
        return {
          type: "paragraph",
          raw: cap[0],
          text: cap[1].charAt(cap[1].length - 1) === "\n" ? cap[1].slice(0, -1) : cap[1]
        };
      }
    }
    text(src) {
      const cap = this.rules.block.text.exec(src);
      if (cap) {
        return {
          type: "text",
          raw: cap[0],
          text: cap[0]
        };
      }
    }
    escape(src) {
      const cap = this.rules.inline.escape.exec(src);
      if (cap) {
        return {
          type: "escape",
          raw: cap[0],
          text: escape(cap[1])
        };
      }
    }
    tag(src, inLink, inRawBlock) {
      const cap = this.rules.inline.tag.exec(src);
      if (cap) {
        if (!inLink && /^<a /i.test(cap[0])) {
          inLink = true;
        } else if (inLink && /^<\/a>/i.test(cap[0])) {
          inLink = false;
        }
        if (!inRawBlock && /^<(pre|code|kbd|script)(\s|>)/i.test(cap[0])) {
          inRawBlock = true;
        } else if (inRawBlock && /^<\/(pre|code|kbd|script)(\s|>)/i.test(cap[0])) {
          inRawBlock = false;
        }
        return {
          type: this.options.sanitize ? "text" : "html",
          raw: cap[0],
          inLink,
          inRawBlock,
          text: this.options.sanitize ? this.options.sanitizer ? this.options.sanitizer(cap[0]) : escape(cap[0]) : cap[0]
        };
      }
    }
    link(src) {
      const cap = this.rules.inline.link.exec(src);
      if (cap) {
        const trimmedUrl = cap[2].trim();
        if (!this.options.pedantic && /^</.test(trimmedUrl)) {
          if (!/>$/.test(trimmedUrl)) {
            return;
          }
          const rtrimSlash = rtrim(trimmedUrl.slice(0, -1), "\\");
          if ((trimmedUrl.length - rtrimSlash.length) % 2 === 0) {
            return;
          }
        } else {
          const lastParenIndex = findClosingBracket(cap[2], "()");
          if (lastParenIndex > -1) {
            const start = cap[0].indexOf("!") === 0 ? 5 : 4;
            const linkLen = start + cap[1].length + lastParenIndex;
            cap[2] = cap[2].substring(0, lastParenIndex);
            cap[0] = cap[0].substring(0, linkLen).trim();
            cap[3] = "";
          }
        }
        let href = cap[2];
        let title = "";
        if (this.options.pedantic) {
          const link2 = /^([^'"]*[^\s])\s+(['"])(.*)\2/.exec(href);
          if (link2) {
            href = link2[1];
            title = link2[3];
          }
        } else {
          title = cap[3] ? cap[3].slice(1, -1) : "";
        }
        href = href.trim();
        if (/^</.test(href)) {
          if (this.options.pedantic && !/>$/.test(trimmedUrl)) {
            href = href.slice(1);
          } else {
            href = href.slice(1, -1);
          }
        }
        return outputLink(cap, {
          href: href ? href.replace(this.rules.inline._escapes, "$1") : href,
          title: title ? title.replace(this.rules.inline._escapes, "$1") : title
        }, cap[0]);
      }
    }
    reflink(src, links) {
      let cap;
      if ((cap = this.rules.inline.reflink.exec(src)) || (cap = this.rules.inline.nolink.exec(src))) {
        let link2 = (cap[2] || cap[1]).replace(/\s+/g, " ");
        link2 = links[link2.toLowerCase()];
        if (!link2 || !link2.href) {
          const text = cap[0].charAt(0);
          return {
            type: "text",
            raw: text,
            text
          };
        }
        return outputLink(cap, link2, cap[0]);
      }
    }
    emStrong(src, maskedSrc, prevChar = "") {
      let match = this.rules.inline.emStrong.lDelim.exec(src);
      if (!match)
        return;
      if (match[3] && prevChar.match(/[\p{L}\p{N}]/u))
        return;
      const nextChar = match[1] || match[2] || "";
      if (!nextChar || nextChar && (prevChar === "" || this.rules.inline.punctuation.exec(prevChar))) {
        const lLength = match[0].length - 1;
        let rDelim, rLength, delimTotal = lLength, midDelimTotal = 0;
        const endReg = match[0][0] === "*" ? this.rules.inline.emStrong.rDelimAst : this.rules.inline.emStrong.rDelimUnd;
        endReg.lastIndex = 0;
        maskedSrc = maskedSrc.slice(-1 * src.length + lLength);
        while ((match = endReg.exec(maskedSrc)) != null) {
          rDelim = match[1] || match[2] || match[3] || match[4] || match[5] || match[6];
          if (!rDelim)
            continue;
          rLength = rDelim.length;
          if (match[3] || match[4]) {
            delimTotal += rLength;
            continue;
          } else if (match[5] || match[6]) {
            if (lLength % 3 && !((lLength + rLength) % 3)) {
              midDelimTotal += rLength;
              continue;
            }
          }
          delimTotal -= rLength;
          if (delimTotal > 0)
            continue;
          if (delimTotal + midDelimTotal - rLength <= 0 && !maskedSrc.slice(endReg.lastIndex).match(endReg)) {
            rLength = Math.min(rLength, rLength + delimTotal + midDelimTotal);
          }
          if (Math.min(lLength, rLength) % 2) {
            return {
              type: "em",
              raw: src.slice(0, lLength + match.index + rLength + 1),
              text: src.slice(1, lLength + match.index + rLength)
            };
          }
          if (Math.min(lLength, rLength) % 2 === 0) {
            return {
              type: "strong",
              raw: src.slice(0, lLength + match.index + rLength + 1),
              text: src.slice(2, lLength + match.index + rLength - 1)
            };
          }
        }
      }
    }
    codespan(src) {
      const cap = this.rules.inline.code.exec(src);
      if (cap) {
        let text = cap[2].replace(/\n/g, " ");
        const hasNonSpaceChars = /[^ ]/.test(text);
        const hasSpaceCharsOnBothEnds = /^ /.test(text) && / $/.test(text);
        if (hasNonSpaceChars && hasSpaceCharsOnBothEnds) {
          text = text.substring(1, text.length - 1);
        }
        text = escape(text, true);
        return {
          type: "codespan",
          raw: cap[0],
          text
        };
      }
    }
    br(src) {
      const cap = this.rules.inline.br.exec(src);
      if (cap) {
        return {
          type: "br",
          raw: cap[0]
        };
      }
    }
    del(src) {
      const cap = this.rules.inline.del.exec(src);
      if (cap) {
        return {
          type: "del",
          raw: cap[0],
          text: cap[2]
        };
      }
    }
    autolink(src, mangle) {
      const cap = this.rules.inline.autolink.exec(src);
      if (cap) {
        let text, href;
        if (cap[2] === "@") {
          text = escape(this.options.mangle ? mangle(cap[1]) : cap[1]);
          href = "mailto:" + text;
        } else {
          text = escape(cap[1]);
          href = text;
        }
        return {
          type: "link",
          raw: cap[0],
          text,
          href,
          tokens: [
            {
              type: "text",
              raw: text,
              text
            }
          ]
        };
      }
    }
    url(src, mangle) {
      let cap;
      if (cap = this.rules.inline.url.exec(src)) {
        let text, href;
        if (cap[2] === "@") {
          text = escape(this.options.mangle ? mangle(cap[0]) : cap[0]);
          href = "mailto:" + text;
        } else {
          let prevCapZero;
          do {
            prevCapZero = cap[0];
            cap[0] = this.rules.inline._backpedal.exec(cap[0])[0];
          } while (prevCapZero !== cap[0]);
          text = escape(cap[0]);
          if (cap[1] === "www.") {
            href = "http://" + text;
          } else {
            href = text;
          }
        }
        return {
          type: "link",
          raw: cap[0],
          text,
          href,
          tokens: [
            {
              type: "text",
              raw: text,
              text
            }
          ]
        };
      }
    }
    inlineText(src, inRawBlock, smartypants) {
      const cap = this.rules.inline.text.exec(src);
      if (cap) {
        let text;
        if (inRawBlock) {
          text = this.options.sanitize ? this.options.sanitizer ? this.options.sanitizer(cap[0]) : escape(cap[0]) : cap[0];
        } else {
          text = escape(this.options.smartypants ? smartypants(cap[0]) : cap[0]);
        }
        return {
          type: "text",
          raw: cap[0],
          text
        };
      }
    }
  };
});

// node_modules/marked/src/rules.js
var require_rules = __commonJS((exports2, module2) => {
  var {
    noopTest,
    edit,
    merge
  } = require_helpers();
  var block = {
    newline: /^(?: *(?:\n|$))+/,
    code: /^( {4}[^\n]+(?:\n(?: *(?:\n|$))*)?)+/,
    fences: /^ {0,3}(`{3,}(?=[^`\n]*\n)|~{3,})([^\n]*)\n(?:|([\s\S]*?)\n)(?: {0,3}\1[~`]* *(?:\n+|$)|$)/,
    hr: /^ {0,3}((?:- *){3,}|(?:_ *){3,}|(?:\* *){3,})(?:\n+|$)/,
    heading: /^ {0,3}(#{1,6})(?=\s|$)(.*)(?:\n+|$)/,
    blockquote: /^( {0,3}> ?(paragraph|[^\n]*)(?:\n|$))+/,
    list: /^( {0,3})(bull) [\s\S]+?(?:hr|def|\n{2,}(?! )(?! {0,3}bull )\n*|\s*$)/,
    html: "^ {0,3}(?:<(script|pre|style)[\\s>][\\s\\S]*?(?:</\\1>[^\\n]*\\n+|$)|comment[^\\n]*(\\n+|$)|<\\?[\\s\\S]*?(?:\\?>\\n*|$)|<![A-Z][\\s\\S]*?(?:>\\n*|$)|<!\\[CDATA\\[[\\s\\S]*?(?:\\]\\]>\\n*|$)|</?(tag)(?: +|\\n|/?>)[\\s\\S]*?(?:\\n{2,}|$)|<(?!script|pre|style)([a-z][\\w-]*)(?:attribute)*? */?>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:\\n{2,}|$)|</(?!script|pre|style)[a-z][\\w-]*\\s*>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:\\n{2,}|$))",
    def: /^ {0,3}\[(label)\]: *\n? *<?([^\s>]+)>?(?:(?: +\n? *| *\n *)(title))? *(?:\n+|$)/,
    nptable: noopTest,
    table: noopTest,
    lheading: /^([^\n]+)\n {0,3}(=+|-+) *(?:\n+|$)/,
    _paragraph: /^([^\n]+(?:\n(?!hr|heading|lheading|blockquote|fences|list|html| +\n)[^\n]+)*)/,
    text: /^[^\n]+/
  };
  block._label = /(?!\s*\])(?:\\[\[\]]|[^\[\]])+/;
  block._title = /(?:"(?:\\"?|[^"\\])*"|'[^'\n]*(?:\n[^'\n]+)*\n?'|\([^()]*\))/;
  block.def = edit(block.def).replace("label", block._label).replace("title", block._title).getRegex();
  block.bullet = /(?:[*+-]|\d{1,9}[.)])/;
  block.item = /^( *)(bull) ?[^\n]*(?:\n(?! *bull ?)[^\n]*)*/;
  block.item = edit(block.item, "gm").replace(/bull/g, block.bullet).getRegex();
  block.listItemStart = edit(/^( *)(bull) */).replace("bull", block.bullet).getRegex();
  block.list = edit(block.list).replace(/bull/g, block.bullet).replace("hr", "\\n+(?=\\1?(?:(?:- *){3,}|(?:_ *){3,}|(?:\\* *){3,})(?:\\n+|$))").replace("def", "\\n+(?=" + block.def.source + ")").getRegex();
  block._tag = "address|article|aside|base|basefont|blockquote|body|caption|center|col|colgroup|dd|details|dialog|dir|div|dl|dt|fieldset|figcaption|figure|footer|form|frame|frameset|h[1-6]|head|header|hr|html|iframe|legend|li|link|main|menu|menuitem|meta|nav|noframes|ol|optgroup|option|p|param|section|source|summary|table|tbody|td|tfoot|th|thead|title|tr|track|ul";
  block._comment = /<!--(?!-?>)[\s\S]*?(?:-->|$)/;
  block.html = edit(block.html, "i").replace("comment", block._comment).replace("tag", block._tag).replace("attribute", / +[a-zA-Z:_][\w.:-]*(?: *= *"[^"\n]*"| *= *'[^'\n]*'| *= *[^\s"'=<>`]+)?/).getRegex();
  block.paragraph = edit(block._paragraph).replace("hr", block.hr).replace("heading", " {0,3}#{1,6} ").replace("|lheading", "").replace("blockquote", " {0,3}>").replace("fences", " {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list", " {0,3}(?:[*+-]|1[.)]) ").replace("html", "</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|!--)").replace("tag", block._tag).getRegex();
  block.blockquote = edit(block.blockquote).replace("paragraph", block.paragraph).getRegex();
  block.normal = merge({}, block);
  block.gfm = merge({}, block.normal, {
    nptable: "^ *([^|\\n ].*\\|.*)\\n {0,3}([-:]+ *\\|[-| :]*)(?:\\n((?:(?!\\n|hr|heading|blockquote|code|fences|list|html).*(?:\\n|$))*)\\n*|$)",
    table: "^ *\\|(.+)\\n {0,3}\\|?( *[-:]+[-| :]*)(?:\\n *((?:(?!\\n|hr|heading|blockquote|code|fences|list|html).*(?:\\n|$))*)\\n*|$)"
  });
  block.gfm.nptable = edit(block.gfm.nptable).replace("hr", block.hr).replace("heading", " {0,3}#{1,6} ").replace("blockquote", " {0,3}>").replace("code", " {4}[^\\n]").replace("fences", " {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list", " {0,3}(?:[*+-]|1[.)]) ").replace("html", "</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|!--)").replace("tag", block._tag).getRegex();
  block.gfm.table = edit(block.gfm.table).replace("hr", block.hr).replace("heading", " {0,3}#{1,6} ").replace("blockquote", " {0,3}>").replace("code", " {4}[^\\n]").replace("fences", " {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list", " {0,3}(?:[*+-]|1[.)]) ").replace("html", "</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|!--)").replace("tag", block._tag).getRegex();
  block.pedantic = merge({}, block.normal, {
    html: edit(`^ *(?:comment *(?:\\n|\\s*$)|<(tag)[\\s\\S]+?</\\1> *(?:\\n{2,}|\\s*$)|<tag(?:"[^"]*"|'[^']*'|\\s[^'"/>\\s]*)*?/?> *(?:\\n{2,}|\\s*$))`).replace("comment", block._comment).replace(/tag/g, "(?!(?:a|em|strong|small|s|cite|q|dfn|abbr|data|time|code|var|samp|kbd|sub|sup|i|b|u|mark|ruby|rt|rp|bdi|bdo|span|br|wbr|ins|del|img)\\b)\\w+(?!:|[^\\w\\s@]*@)\\b").getRegex(),
    def: /^ *\[([^\]]+)\]: *<?([^\s>]+)>?(?: +(["(][^\n]+[")]))? *(?:\n+|$)/,
    heading: /^(#{1,6})(.*)(?:\n+|$)/,
    fences: noopTest,
    paragraph: edit(block.normal._paragraph).replace("hr", block.hr).replace("heading", " *#{1,6} *[^\n]").replace("lheading", block.lheading).replace("blockquote", " {0,3}>").replace("|fences", "").replace("|list", "").replace("|html", "").getRegex()
  });
  var inline = {
    escape: /^\\([!"#$%&'()*+,\-./:;<=>?@\[\]\\^_`{|}~])/,
    autolink: /^<(scheme:[^\s\x00-\x1f<>]*|email)>/,
    url: noopTest,
    tag: "^comment|^</[a-zA-Z][\\w:-]*\\s*>|^<[a-zA-Z][\\w-]*(?:attribute)*?\\s*/?>|^<\\?[\\s\\S]*?\\?>|^<![a-zA-Z]+\\s[\\s\\S]*?>|^<!\\[CDATA\\[[\\s\\S]*?\\]\\]>",
    link: /^!?\[(label)\]\(\s*(href)(?:\s+(title))?\s*\)/,
    reflink: /^!?\[(label)\]\[(?!\s*\])((?:\\[\[\]]?|[^\[\]\\])+)\]/,
    nolink: /^!?\[(?!\s*\])((?:\[[^\[\]]*\]|\\[\[\]]|[^\[\]])*)\](?:\[\])?/,
    reflinkSearch: "reflink|nolink(?!\\()",
    emStrong: {
      lDelim: /^(?:\*+(?:([punct_])|[^\s*]))|^_+(?:([punct*])|([^\s_]))/,
      rDelimAst: /\_\_[^_]*?\*[^_]*?\_\_|[punct_](\*+)(?=[\s]|$)|[^punct*_\s](\*+)(?=[punct_\s]|$)|[punct_\s](\*+)(?=[^punct*_\s])|[\s](\*+)(?=[punct_])|[punct_](\*+)(?=[punct_])|[^punct*_\s](\*+)(?=[^punct*_\s])/,
      rDelimUnd: /\*\*[^*]*?\_[^*]*?\*\*|[punct*](\_+)(?=[\s]|$)|[^punct*_\s](\_+)(?=[punct*\s]|$)|[punct*\s](\_+)(?=[^punct*_\s])|[\s](\_+)(?=[punct*])|[punct*](\_+)(?=[punct*])/
    },
    code: /^(`+)([^`]|[^`][\s\S]*?[^`])\1(?!`)/,
    br: /^( {2,}|\\)\n(?!\s*$)/,
    del: noopTest,
    text: /^(`+|[^`])(?:(?= {2,}\n)|[\s\S]*?(?:(?=[\\<!\[`*_]|\b_|$)|[^ ](?= {2,}\n)))/,
    punctuation: /^([\spunctuation])/
  };
  inline._punctuation = "!\"#$%&'()+\\-.,/:;<=>?@\\[\\]`^{|}~";
  inline.punctuation = edit(inline.punctuation).replace(/punctuation/g, inline._punctuation).getRegex();
  inline.blockSkip = /\[[^\]]*?\]\([^\)]*?\)|`[^`]*?`|<[^>]*?>/g;
  inline.escapedEmSt = /\\\*|\\_/g;
  inline._comment = edit(block._comment).replace("(?:-->|$)", "-->").getRegex();
  inline.emStrong.lDelim = edit(inline.emStrong.lDelim).replace(/punct/g, inline._punctuation).getRegex();
  inline.emStrong.rDelimAst = edit(inline.emStrong.rDelimAst, "g").replace(/punct/g, inline._punctuation).getRegex();
  inline.emStrong.rDelimUnd = edit(inline.emStrong.rDelimUnd, "g").replace(/punct/g, inline._punctuation).getRegex();
  inline._escapes = /\\([!"#$%&'()*+,\-./:;<=>?@\[\]\\^_`{|}~])/g;
  inline._scheme = /[a-zA-Z][a-zA-Z0-9+.-]{1,31}/;
  inline._email = /[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+(@)[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+(?![-_])/;
  inline.autolink = edit(inline.autolink).replace("scheme", inline._scheme).replace("email", inline._email).getRegex();
  inline._attribute = /\s+[a-zA-Z:_][\w.:-]*(?:\s*=\s*"[^"]*"|\s*=\s*'[^']*'|\s*=\s*[^\s"'=<>`]+)?/;
  inline.tag = edit(inline.tag).replace("comment", inline._comment).replace("attribute", inline._attribute).getRegex();
  inline._label = /(?:\[(?:\\.|[^\[\]\\])*\]|\\.|`[^`]*`|[^\[\]\\`])*?/;
  inline._href = /<(?:\\.|[^\n<>\\])+>|[^\s\x00-\x1f]*/;
  inline._title = /"(?:\\"?|[^"\\])*"|'(?:\\'?|[^'\\])*'|\((?:\\\)?|[^)\\])*\)/;
  inline.link = edit(inline.link).replace("label", inline._label).replace("href", inline._href).replace("title", inline._title).getRegex();
  inline.reflink = edit(inline.reflink).replace("label", inline._label).getRegex();
  inline.reflinkSearch = edit(inline.reflinkSearch, "g").replace("reflink", inline.reflink).replace("nolink", inline.nolink).getRegex();
  inline.normal = merge({}, inline);
  inline.pedantic = merge({}, inline.normal, {
    strong: {
      start: /^__|\*\*/,
      middle: /^__(?=\S)([\s\S]*?\S)__(?!_)|^\*\*(?=\S)([\s\S]*?\S)\*\*(?!\*)/,
      endAst: /\*\*(?!\*)/g,
      endUnd: /__(?!_)/g
    },
    em: {
      start: /^_|\*/,
      middle: /^()\*(?=\S)([\s\S]*?\S)\*(?!\*)|^_(?=\S)([\s\S]*?\S)_(?!_)/,
      endAst: /\*(?!\*)/g,
      endUnd: /_(?!_)/g
    },
    link: edit(/^!?\[(label)\]\((.*?)\)/).replace("label", inline._label).getRegex(),
    reflink: edit(/^!?\[(label)\]\s*\[([^\]]*)\]/).replace("label", inline._label).getRegex()
  });
  inline.gfm = merge({}, inline.normal, {
    escape: edit(inline.escape).replace("])", "~|])").getRegex(),
    _extended_email: /[A-Za-z0-9._+-]+(@)[a-zA-Z0-9-_]+(?:\.[a-zA-Z0-9-_]*[a-zA-Z0-9])+(?![-_])/,
    url: /^((?:ftp|https?):\/\/|www\.)(?:[a-zA-Z0-9\-]+\.?)+[^\s<]*|^email/,
    _backpedal: /(?:[^?!.,:;*_~()&]+|\([^)]*\)|&(?![a-zA-Z0-9]+;$)|[?!.,:;*_~)]+(?!$))+/,
    del: /^(~~?)(?=[^\s~])([\s\S]*?[^\s~])\1(?=[^~]|$)/,
    text: /^([`~]+|[^`~])(?:(?= {2,}\n)|[\s\S]*?(?:(?=[\\<!\[`*~_]|\b_|https?:\/\/|ftp:\/\/|www\.|$)|[^ ](?= {2,}\n)|[^a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-](?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@))|(?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@))/
  });
  inline.gfm.url = edit(inline.gfm.url, "i").replace("email", inline.gfm._extended_email).getRegex();
  inline.breaks = merge({}, inline.gfm, {
    br: edit(inline.br).replace("{2,}", "*").getRegex(),
    text: edit(inline.gfm.text).replace("\\b_", "\\b_| {2,}\\n").replace(/\{2,\}/g, "*").getRegex()
  });
  module2.exports = {
    block,
    inline
  };
});

// node_modules/marked/src/Lexer.js
var require_Lexer = __commonJS((exports2, module2) => {
  var Tokenizer = require_Tokenizer();
  var {defaults} = require_defaults();
  var {block, inline} = require_rules();
  var {repeatString} = require_helpers();
  function smartypants(text) {
    return text.replace(/---/g, "\u2014").replace(/--/g, "\u2013").replace(/(^|[-\u2014/(\[{"\s])'/g, "$1\u2018").replace(/'/g, "\u2019").replace(/(^|[-\u2014/(\[{\u2018\s])"/g, "$1\u201C").replace(/"/g, "\u201D").replace(/\.{3}/g, "\u2026");
  }
  function mangle(text) {
    let out = "", i, ch;
    const l = text.length;
    for (i = 0; i < l; i++) {
      ch = text.charCodeAt(i);
      if (Math.random() > 0.5) {
        ch = "x" + ch.toString(16);
      }
      out += "&#" + ch + ";";
    }
    return out;
  }
  module2.exports = class Lexer {
    constructor(options) {
      this.tokens = [];
      this.tokens.links = Object.create(null);
      this.options = options || defaults;
      this.options.tokenizer = this.options.tokenizer || new Tokenizer();
      this.tokenizer = this.options.tokenizer;
      this.tokenizer.options = this.options;
      const rules = {
        block: block.normal,
        inline: inline.normal
      };
      if (this.options.pedantic) {
        rules.block = block.pedantic;
        rules.inline = inline.pedantic;
      } else if (this.options.gfm) {
        rules.block = block.gfm;
        if (this.options.breaks) {
          rules.inline = inline.breaks;
        } else {
          rules.inline = inline.gfm;
        }
      }
      this.tokenizer.rules = rules;
    }
    static get rules() {
      return {
        block,
        inline
      };
    }
    static lex(src, options) {
      const lexer = new Lexer(options);
      return lexer.lex(src);
    }
    static lexInline(src, options) {
      const lexer = new Lexer(options);
      return lexer.inlineTokens(src);
    }
    lex(src) {
      src = src.replace(/\r\n|\r/g, "\n").replace(/\t/g, "    ");
      this.blockTokens(src, this.tokens, true);
      this.inline(this.tokens);
      return this.tokens;
    }
    blockTokens(src, tokens = [], top = true) {
      if (this.options.pedantic) {
        src = src.replace(/^ +$/gm, "");
      }
      let token, i, l, lastToken;
      while (src) {
        if (token = this.tokenizer.space(src)) {
          src = src.substring(token.raw.length);
          if (token.type) {
            tokens.push(token);
          }
          continue;
        }
        if (token = this.tokenizer.code(src)) {
          src = src.substring(token.raw.length);
          lastToken = tokens[tokens.length - 1];
          if (lastToken && lastToken.type === "paragraph") {
            lastToken.raw += "\n" + token.raw;
            lastToken.text += "\n" + token.text;
          } else {
            tokens.push(token);
          }
          continue;
        }
        if (token = this.tokenizer.fences(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.heading(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.nptable(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.hr(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.blockquote(src)) {
          src = src.substring(token.raw.length);
          token.tokens = this.blockTokens(token.text, [], top);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.list(src)) {
          src = src.substring(token.raw.length);
          l = token.items.length;
          for (i = 0; i < l; i++) {
            token.items[i].tokens = this.blockTokens(token.items[i].text, [], false);
          }
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.html(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (top && (token = this.tokenizer.def(src))) {
          src = src.substring(token.raw.length);
          if (!this.tokens.links[token.tag]) {
            this.tokens.links[token.tag] = {
              href: token.href,
              title: token.title
            };
          }
          continue;
        }
        if (token = this.tokenizer.table(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.lheading(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (top && (token = this.tokenizer.paragraph(src))) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.text(src)) {
          src = src.substring(token.raw.length);
          lastToken = tokens[tokens.length - 1];
          if (lastToken && lastToken.type === "text") {
            lastToken.raw += "\n" + token.raw;
            lastToken.text += "\n" + token.text;
          } else {
            tokens.push(token);
          }
          continue;
        }
        if (src) {
          const errMsg = "Infinite loop on byte: " + src.charCodeAt(0);
          if (this.options.silent) {
            console.error(errMsg);
            break;
          } else {
            throw new Error(errMsg);
          }
        }
      }
      return tokens;
    }
    inline(tokens) {
      let i, j, k, l2, row, token;
      const l = tokens.length;
      for (i = 0; i < l; i++) {
        token = tokens[i];
        switch (token.type) {
          case "paragraph":
          case "text":
          case "heading": {
            token.tokens = [];
            this.inlineTokens(token.text, token.tokens);
            break;
          }
          case "table": {
            token.tokens = {
              header: [],
              cells: []
            };
            l2 = token.header.length;
            for (j = 0; j < l2; j++) {
              token.tokens.header[j] = [];
              this.inlineTokens(token.header[j], token.tokens.header[j]);
            }
            l2 = token.cells.length;
            for (j = 0; j < l2; j++) {
              row = token.cells[j];
              token.tokens.cells[j] = [];
              for (k = 0; k < row.length; k++) {
                token.tokens.cells[j][k] = [];
                this.inlineTokens(row[k], token.tokens.cells[j][k]);
              }
            }
            break;
          }
          case "blockquote": {
            this.inline(token.tokens);
            break;
          }
          case "list": {
            l2 = token.items.length;
            for (j = 0; j < l2; j++) {
              this.inline(token.items[j].tokens);
            }
            break;
          }
          default: {
          }
        }
      }
      return tokens;
    }
    inlineTokens(src, tokens = [], inLink = false, inRawBlock = false) {
      let token, lastToken;
      let maskedSrc = src;
      let match;
      let keepPrevChar, prevChar;
      if (this.tokens.links) {
        const links = Object.keys(this.tokens.links);
        if (links.length > 0) {
          while ((match = this.tokenizer.rules.inline.reflinkSearch.exec(maskedSrc)) != null) {
            if (links.includes(match[0].slice(match[0].lastIndexOf("[") + 1, -1))) {
              maskedSrc = maskedSrc.slice(0, match.index) + "[" + repeatString("a", match[0].length - 2) + "]" + maskedSrc.slice(this.tokenizer.rules.inline.reflinkSearch.lastIndex);
            }
          }
        }
      }
      while ((match = this.tokenizer.rules.inline.blockSkip.exec(maskedSrc)) != null) {
        maskedSrc = maskedSrc.slice(0, match.index) + "[" + repeatString("a", match[0].length - 2) + "]" + maskedSrc.slice(this.tokenizer.rules.inline.blockSkip.lastIndex);
      }
      while ((match = this.tokenizer.rules.inline.escapedEmSt.exec(maskedSrc)) != null) {
        maskedSrc = maskedSrc.slice(0, match.index) + "++" + maskedSrc.slice(this.tokenizer.rules.inline.escapedEmSt.lastIndex);
      }
      while (src) {
        if (!keepPrevChar) {
          prevChar = "";
        }
        keepPrevChar = false;
        if (token = this.tokenizer.escape(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.tag(src, inLink, inRawBlock)) {
          src = src.substring(token.raw.length);
          inLink = token.inLink;
          inRawBlock = token.inRawBlock;
          const lastToken2 = tokens[tokens.length - 1];
          if (lastToken2 && token.type === "text" && lastToken2.type === "text") {
            lastToken2.raw += token.raw;
            lastToken2.text += token.text;
          } else {
            tokens.push(token);
          }
          continue;
        }
        if (token = this.tokenizer.link(src)) {
          src = src.substring(token.raw.length);
          if (token.type === "link") {
            token.tokens = this.inlineTokens(token.text, [], true, inRawBlock);
          }
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.reflink(src, this.tokens.links)) {
          src = src.substring(token.raw.length);
          const lastToken2 = tokens[tokens.length - 1];
          if (token.type === "link") {
            token.tokens = this.inlineTokens(token.text, [], true, inRawBlock);
            tokens.push(token);
          } else if (lastToken2 && token.type === "text" && lastToken2.type === "text") {
            lastToken2.raw += token.raw;
            lastToken2.text += token.text;
          } else {
            tokens.push(token);
          }
          continue;
        }
        if (token = this.tokenizer.emStrong(src, maskedSrc, prevChar)) {
          src = src.substring(token.raw.length);
          token.tokens = this.inlineTokens(token.text, [], inLink, inRawBlock);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.codespan(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.br(src)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.del(src)) {
          src = src.substring(token.raw.length);
          token.tokens = this.inlineTokens(token.text, [], inLink, inRawBlock);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.autolink(src, mangle)) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (!inLink && (token = this.tokenizer.url(src, mangle))) {
          src = src.substring(token.raw.length);
          tokens.push(token);
          continue;
        }
        if (token = this.tokenizer.inlineText(src, inRawBlock, smartypants)) {
          src = src.substring(token.raw.length);
          if (token.raw.slice(-1) !== "_") {
            prevChar = token.raw.slice(-1);
          }
          keepPrevChar = true;
          lastToken = tokens[tokens.length - 1];
          if (lastToken && lastToken.type === "text") {
            lastToken.raw += token.raw;
            lastToken.text += token.text;
          } else {
            tokens.push(token);
          }
          continue;
        }
        if (src) {
          const errMsg = "Infinite loop on byte: " + src.charCodeAt(0);
          if (this.options.silent) {
            console.error(errMsg);
            break;
          } else {
            throw new Error(errMsg);
          }
        }
      }
      return tokens;
    }
  };
});

// node_modules/marked/src/Renderer.js
var require_Renderer = __commonJS((exports2, module2) => {
  var {defaults} = require_defaults();
  var {
    cleanUrl,
    escape
  } = require_helpers();
  module2.exports = class Renderer {
    constructor(options) {
      this.options = options || defaults;
    }
    code(code, infostring, escaped) {
      const lang = (infostring || "").match(/\S*/)[0];
      if (this.options.highlight) {
        const out = this.options.highlight(code, lang);
        if (out != null && out !== code) {
          escaped = true;
          code = out;
        }
      }
      code = code.replace(/\n$/, "") + "\n";
      if (!lang) {
        return "<pre><code>" + (escaped ? code : escape(code, true)) + "</code></pre>\n";
      }
      return '<pre><code class="' + this.options.langPrefix + escape(lang, true) + '">' + (escaped ? code : escape(code, true)) + "</code></pre>\n";
    }
    blockquote(quote) {
      return "<blockquote>\n" + quote + "</blockquote>\n";
    }
    html(html) {
      return html;
    }
    heading(text, level, raw, slugger) {
      if (this.options.headerIds) {
        return "<h" + level + ' id="' + this.options.headerPrefix + slugger.slug(raw) + '">' + text + "</h" + level + ">\n";
      }
      return "<h" + level + ">" + text + "</h" + level + ">\n";
    }
    hr() {
      return this.options.xhtml ? "<hr/>\n" : "<hr>\n";
    }
    list(body, ordered, start) {
      const type = ordered ? "ol" : "ul", startatt = ordered && start !== 1 ? ' start="' + start + '"' : "";
      return "<" + type + startatt + ">\n" + body + "</" + type + ">\n";
    }
    listitem(text) {
      return "<li>" + text + "</li>\n";
    }
    checkbox(checked) {
      return "<input " + (checked ? 'checked="" ' : "") + 'disabled="" type="checkbox"' + (this.options.xhtml ? " /" : "") + "> ";
    }
    paragraph(text) {
      return "<p>" + text + "</p>\n";
    }
    table(header, body) {
      if (body)
        body = "<tbody>" + body + "</tbody>";
      return "<table>\n<thead>\n" + header + "</thead>\n" + body + "</table>\n";
    }
    tablerow(content) {
      return "<tr>\n" + content + "</tr>\n";
    }
    tablecell(content, flags) {
      const type = flags.header ? "th" : "td";
      const tag = flags.align ? "<" + type + ' align="' + flags.align + '">' : "<" + type + ">";
      return tag + content + "</" + type + ">\n";
    }
    strong(text) {
      return "<strong>" + text + "</strong>";
    }
    em(text) {
      return "<em>" + text + "</em>";
    }
    codespan(text) {
      return "<code>" + text + "</code>";
    }
    br() {
      return this.options.xhtml ? "<br/>" : "<br>";
    }
    del(text) {
      return "<del>" + text + "</del>";
    }
    link(href, title, text) {
      href = cleanUrl(this.options.sanitize, this.options.baseUrl, href);
      if (href === null) {
        return text;
      }
      let out = '<a href="' + escape(href) + '"';
      if (title) {
        out += ' title="' + title + '"';
      }
      out += ">" + text + "</a>";
      return out;
    }
    image(href, title, text) {
      href = cleanUrl(this.options.sanitize, this.options.baseUrl, href);
      if (href === null) {
        return text;
      }
      let out = '<img src="' + href + '" alt="' + text + '"';
      if (title) {
        out += ' title="' + title + '"';
      }
      out += this.options.xhtml ? "/>" : ">";
      return out;
    }
    text(text) {
      return text;
    }
  };
});

// node_modules/marked/src/TextRenderer.js
var require_TextRenderer = __commonJS((exports2, module2) => {
  module2.exports = class TextRenderer {
    strong(text) {
      return text;
    }
    em(text) {
      return text;
    }
    codespan(text) {
      return text;
    }
    del(text) {
      return text;
    }
    html(text) {
      return text;
    }
    text(text) {
      return text;
    }
    link(href, title, text) {
      return "" + text;
    }
    image(href, title, text) {
      return "" + text;
    }
    br() {
      return "";
    }
  };
});

// node_modules/marked/src/Slugger.js
var require_Slugger = __commonJS((exports2, module2) => {
  module2.exports = class Slugger {
    constructor() {
      this.seen = {};
    }
    serialize(value) {
      return value.toLowerCase().trim().replace(/<[!\/a-z].*?>/ig, "").replace(/[\u2000-\u206F\u2E00-\u2E7F\\'!"#$%&()*+,./:;<=>?@[\]^`{|}~]/g, "").replace(/\s/g, "-");
    }
    getNextSafeSlug(originalSlug, isDryRun) {
      let slug = originalSlug;
      let occurenceAccumulator = 0;
      if (this.seen.hasOwnProperty(slug)) {
        occurenceAccumulator = this.seen[originalSlug];
        do {
          occurenceAccumulator++;
          slug = originalSlug + "-" + occurenceAccumulator;
        } while (this.seen.hasOwnProperty(slug));
      }
      if (!isDryRun) {
        this.seen[originalSlug] = occurenceAccumulator;
        this.seen[slug] = 0;
      }
      return slug;
    }
    slug(value, options = {}) {
      const slug = this.serialize(value);
      return this.getNextSafeSlug(slug, options.dryrun);
    }
  };
});

// node_modules/marked/src/Parser.js
var require_Parser = __commonJS((exports2, module2) => {
  var Renderer = require_Renderer();
  var TextRenderer = require_TextRenderer();
  var Slugger = require_Slugger();
  var {defaults} = require_defaults();
  var {
    unescape
  } = require_helpers();
  module2.exports = class Parser {
    constructor(options) {
      this.options = options || defaults;
      this.options.renderer = this.options.renderer || new Renderer();
      this.renderer = this.options.renderer;
      this.renderer.options = this.options;
      this.textRenderer = new TextRenderer();
      this.slugger = new Slugger();
    }
    static parse(tokens, options) {
      const parser = new Parser(options);
      return parser.parse(tokens);
    }
    static parseInline(tokens, options) {
      const parser = new Parser(options);
      return parser.parseInline(tokens);
    }
    parse(tokens, top = true) {
      let out = "", i, j, k, l2, l3, row, cell, header, body, token, ordered, start, loose, itemBody, item, checked, task, checkbox;
      const l = tokens.length;
      for (i = 0; i < l; i++) {
        token = tokens[i];
        switch (token.type) {
          case "space": {
            continue;
          }
          case "hr": {
            out += this.renderer.hr();
            continue;
          }
          case "heading": {
            out += this.renderer.heading(this.parseInline(token.tokens), token.depth, unescape(this.parseInline(token.tokens, this.textRenderer)), this.slugger);
            continue;
          }
          case "code": {
            out += this.renderer.code(token.text, token.lang, token.escaped);
            continue;
          }
          case "table": {
            header = "";
            cell = "";
            l2 = token.header.length;
            for (j = 0; j < l2; j++) {
              cell += this.renderer.tablecell(this.parseInline(token.tokens.header[j]), {header: true, align: token.align[j]});
            }
            header += this.renderer.tablerow(cell);
            body = "";
            l2 = token.cells.length;
            for (j = 0; j < l2; j++) {
              row = token.tokens.cells[j];
              cell = "";
              l3 = row.length;
              for (k = 0; k < l3; k++) {
                cell += this.renderer.tablecell(this.parseInline(row[k]), {header: false, align: token.align[k]});
              }
              body += this.renderer.tablerow(cell);
            }
            out += this.renderer.table(header, body);
            continue;
          }
          case "blockquote": {
            body = this.parse(token.tokens);
            out += this.renderer.blockquote(body);
            continue;
          }
          case "list": {
            ordered = token.ordered;
            start = token.start;
            loose = token.loose;
            l2 = token.items.length;
            body = "";
            for (j = 0; j < l2; j++) {
              item = token.items[j];
              checked = item.checked;
              task = item.task;
              itemBody = "";
              if (item.task) {
                checkbox = this.renderer.checkbox(checked);
                if (loose) {
                  if (item.tokens.length > 0 && item.tokens[0].type === "text") {
                    item.tokens[0].text = checkbox + " " + item.tokens[0].text;
                    if (item.tokens[0].tokens && item.tokens[0].tokens.length > 0 && item.tokens[0].tokens[0].type === "text") {
                      item.tokens[0].tokens[0].text = checkbox + " " + item.tokens[0].tokens[0].text;
                    }
                  } else {
                    item.tokens.unshift({
                      type: "text",
                      text: checkbox
                    });
                  }
                } else {
                  itemBody += checkbox;
                }
              }
              itemBody += this.parse(item.tokens, loose);
              body += this.renderer.listitem(itemBody, task, checked);
            }
            out += this.renderer.list(body, ordered, start);
            continue;
          }
          case "html": {
            out += this.renderer.html(token.text);
            continue;
          }
          case "paragraph": {
            out += this.renderer.paragraph(this.parseInline(token.tokens));
            continue;
          }
          case "text": {
            body = token.tokens ? this.parseInline(token.tokens) : token.text;
            while (i + 1 < l && tokens[i + 1].type === "text") {
              token = tokens[++i];
              body += "\n" + (token.tokens ? this.parseInline(token.tokens) : token.text);
            }
            out += top ? this.renderer.paragraph(body) : body;
            continue;
          }
          default: {
            const errMsg = 'Token with "' + token.type + '" type was not found.';
            if (this.options.silent) {
              console.error(errMsg);
              return;
            } else {
              throw new Error(errMsg);
            }
          }
        }
      }
      return out;
    }
    parseInline(tokens, renderer) {
      renderer = renderer || this.renderer;
      let out = "", i, token;
      const l = tokens.length;
      for (i = 0; i < l; i++) {
        token = tokens[i];
        switch (token.type) {
          case "escape": {
            out += renderer.text(token.text);
            break;
          }
          case "html": {
            out += renderer.html(token.text);
            break;
          }
          case "link": {
            out += renderer.link(token.href, token.title, this.parseInline(token.tokens, renderer));
            break;
          }
          case "image": {
            out += renderer.image(token.href, token.title, token.text);
            break;
          }
          case "strong": {
            out += renderer.strong(this.parseInline(token.tokens, renderer));
            break;
          }
          case "em": {
            out += renderer.em(this.parseInline(token.tokens, renderer));
            break;
          }
          case "codespan": {
            out += renderer.codespan(token.text);
            break;
          }
          case "br": {
            out += renderer.br();
            break;
          }
          case "del": {
            out += renderer.del(this.parseInline(token.tokens, renderer));
            break;
          }
          case "text": {
            out += renderer.text(token.text);
            break;
          }
          default: {
            const errMsg = 'Token with "' + token.type + '" type was not found.';
            if (this.options.silent) {
              console.error(errMsg);
              return;
            } else {
              throw new Error(errMsg);
            }
          }
        }
      }
      return out;
    }
  };
});

// node_modules/marked/src/marked.js
var require_marked = __commonJS((exports2, module2) => {
  var Lexer = require_Lexer();
  var Parser = require_Parser();
  var Tokenizer = require_Tokenizer();
  var Renderer = require_Renderer();
  var TextRenderer = require_TextRenderer();
  var Slugger = require_Slugger();
  var {
    merge,
    checkSanitizeDeprecation,
    escape
  } = require_helpers();
  var {
    getDefaults,
    changeDefaults,
    defaults
  } = require_defaults();
  function marked2(src, opt, callback) {
    if (typeof src === "undefined" || src === null) {
      throw new Error("marked(): input parameter is undefined or null");
    }
    if (typeof src !== "string") {
      throw new Error("marked(): input parameter is of type " + Object.prototype.toString.call(src) + ", string expected");
    }
    if (typeof opt === "function") {
      callback = opt;
      opt = null;
    }
    opt = merge({}, marked2.defaults, opt || {});
    checkSanitizeDeprecation(opt);
    if (callback) {
      const highlight = opt.highlight;
      let tokens;
      try {
        tokens = Lexer.lex(src, opt);
      } catch (e) {
        return callback(e);
      }
      const done = function(err) {
        let out;
        if (!err) {
          try {
            out = Parser.parse(tokens, opt);
          } catch (e) {
            err = e;
          }
        }
        opt.highlight = highlight;
        return err ? callback(err) : callback(null, out);
      };
      if (!highlight || highlight.length < 3) {
        return done();
      }
      delete opt.highlight;
      if (!tokens.length)
        return done();
      let pending = 0;
      marked2.walkTokens(tokens, function(token) {
        if (token.type === "code") {
          pending++;
          setTimeout(() => {
            highlight(token.text, token.lang, function(err, code) {
              if (err) {
                return done(err);
              }
              if (code != null && code !== token.text) {
                token.text = code;
                token.escaped = true;
              }
              pending--;
              if (pending === 0) {
                done();
              }
            });
          }, 0);
        }
      });
      if (pending === 0) {
        done();
      }
      return;
    }
    try {
      const tokens = Lexer.lex(src, opt);
      if (opt.walkTokens) {
        marked2.walkTokens(tokens, opt.walkTokens);
      }
      return Parser.parse(tokens, opt);
    } catch (e) {
      e.message += "\nPlease report this to https://github.com/markedjs/marked.";
      if (opt.silent) {
        return "<p>An error occurred:</p><pre>" + escape(e.message + "", true) + "</pre>";
      }
      throw e;
    }
  }
  marked2.options = marked2.setOptions = function(opt) {
    merge(marked2.defaults, opt);
    changeDefaults(marked2.defaults);
    return marked2;
  };
  marked2.getDefaults = getDefaults;
  marked2.defaults = defaults;
  marked2.use = function(extension) {
    const opts = merge({}, extension);
    if (extension.renderer) {
      const renderer = marked2.defaults.renderer || new Renderer();
      for (const prop in extension.renderer) {
        const prevRenderer = renderer[prop];
        renderer[prop] = (...args) => {
          let ret = extension.renderer[prop].apply(renderer, args);
          if (ret === false) {
            ret = prevRenderer.apply(renderer, args);
          }
          return ret;
        };
      }
      opts.renderer = renderer;
    }
    if (extension.tokenizer) {
      const tokenizer = marked2.defaults.tokenizer || new Tokenizer();
      for (const prop in extension.tokenizer) {
        const prevTokenizer = tokenizer[prop];
        tokenizer[prop] = (...args) => {
          let ret = extension.tokenizer[prop].apply(tokenizer, args);
          if (ret === false) {
            ret = prevTokenizer.apply(tokenizer, args);
          }
          return ret;
        };
      }
      opts.tokenizer = tokenizer;
    }
    if (extension.walkTokens) {
      const walkTokens = marked2.defaults.walkTokens;
      opts.walkTokens = (token) => {
        extension.walkTokens(token);
        if (walkTokens) {
          walkTokens(token);
        }
      };
    }
    marked2.setOptions(opts);
  };
  marked2.walkTokens = function(tokens, callback) {
    for (const token of tokens) {
      callback(token);
      switch (token.type) {
        case "table": {
          for (const cell of token.tokens.header) {
            marked2.walkTokens(cell, callback);
          }
          for (const row of token.tokens.cells) {
            for (const cell of row) {
              marked2.walkTokens(cell, callback);
            }
          }
          break;
        }
        case "list": {
          marked2.walkTokens(token.items, callback);
          break;
        }
        default: {
          if (token.tokens) {
            marked2.walkTokens(token.tokens, callback);
          }
        }
      }
    }
  };
  marked2.parseInline = function(src, opt) {
    if (typeof src === "undefined" || src === null) {
      throw new Error("marked.parseInline(): input parameter is undefined or null");
    }
    if (typeof src !== "string") {
      throw new Error("marked.parseInline(): input parameter is of type " + Object.prototype.toString.call(src) + ", string expected");
    }
    opt = merge({}, marked2.defaults, opt || {});
    checkSanitizeDeprecation(opt);
    try {
      const tokens = Lexer.lexInline(src, opt);
      if (opt.walkTokens) {
        marked2.walkTokens(tokens, opt.walkTokens);
      }
      return Parser.parseInline(tokens, opt);
    } catch (e) {
      e.message += "\nPlease report this to https://github.com/markedjs/marked.";
      if (opt.silent) {
        return "<p>An error occurred:</p><pre>" + escape(e.message + "", true) + "</pre>";
      }
      throw e;
    }
  };
  marked2.Parser = Parser;
  marked2.parser = Parser.parse;
  marked2.Renderer = Renderer;
  marked2.TextRenderer = TextRenderer;
  marked2.Lexer = Lexer;
  marked2.lexer = Lexer.lex;
  marked2.Tokenizer = Tokenizer;
  marked2.Slugger = Slugger;
  marked2.parse = marked2;
  module2.exports = marked2;
});

// node_modules/dompurify/dist/purify.cjs.js
var require_purify_cjs = __commonJS((exports2, module2) => {
  /*! @license DOMPurify | (c) Cure53 and other contributors | Released under the Apache license 2.0 and Mozilla Public License 2.0 | github.com/cure53/DOMPurify/blob/2.2.2/LICENSE */
  "use strict";
  function _toConsumableArray(arr) {
    if (Array.isArray(arr)) {
      for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) {
        arr2[i] = arr[i];
      }
      return arr2;
    } else {
      return Array.from(arr);
    }
  }
  var hasOwnProperty = Object.hasOwnProperty;
  var setPrototypeOf = Object.setPrototypeOf;
  var isFrozen = Object.isFrozen;
  var getPrototypeOf = Object.getPrototypeOf;
  var getOwnPropertyDescriptor = Object.getOwnPropertyDescriptor;
  var freeze = Object.freeze;
  var seal = Object.seal;
  var create = Object.create;
  var _ref = typeof Reflect !== "undefined" && Reflect;
  var apply = _ref.apply;
  var construct = _ref.construct;
  if (!apply) {
    apply = function apply2(fun, thisValue, args) {
      return fun.apply(thisValue, args);
    };
  }
  if (!freeze) {
    freeze = function freeze2(x) {
      return x;
    };
  }
  if (!seal) {
    seal = function seal2(x) {
      return x;
    };
  }
  if (!construct) {
    construct = function construct2(Func, args) {
      return new (Function.prototype.bind.apply(Func, [null].concat(_toConsumableArray(args))))();
    };
  }
  var arrayForEach = unapply(Array.prototype.forEach);
  var arrayPop = unapply(Array.prototype.pop);
  var arrayPush = unapply(Array.prototype.push);
  var stringToLowerCase = unapply(String.prototype.toLowerCase);
  var stringMatch = unapply(String.prototype.match);
  var stringReplace = unapply(String.prototype.replace);
  var stringIndexOf = unapply(String.prototype.indexOf);
  var stringTrim = unapply(String.prototype.trim);
  var regExpTest = unapply(RegExp.prototype.test);
  var typeErrorCreate = unconstruct(TypeError);
  function unapply(func) {
    return function(thisArg) {
      for (var _len = arguments.length, args = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
        args[_key - 1] = arguments[_key];
      }
      return apply(func, thisArg, args);
    };
  }
  function unconstruct(func) {
    return function() {
      for (var _len2 = arguments.length, args = Array(_len2), _key2 = 0; _key2 < _len2; _key2++) {
        args[_key2] = arguments[_key2];
      }
      return construct(func, args);
    };
  }
  function addToSet(set3, array) {
    if (setPrototypeOf) {
      setPrototypeOf(set3, null);
    }
    var l = array.length;
    while (l--) {
      var element = array[l];
      if (typeof element === "string") {
        var lcElement = stringToLowerCase(element);
        if (lcElement !== element) {
          if (!isFrozen(array)) {
            array[l] = lcElement;
          }
          element = lcElement;
        }
      }
      set3[element] = true;
    }
    return set3;
  }
  function clone(object) {
    var newObject = create(null);
    var property = void 0;
    for (property in object) {
      if (apply(hasOwnProperty, object, [property])) {
        newObject[property] = object[property];
      }
    }
    return newObject;
  }
  function lookupGetter(object, prop) {
    while (object !== null) {
      var desc = getOwnPropertyDescriptor(object, prop);
      if (desc) {
        if (desc.get) {
          return unapply(desc.get);
        }
        if (typeof desc.value === "function") {
          return unapply(desc.value);
        }
      }
      object = getPrototypeOf(object);
    }
    return null;
  }
  var html = freeze(["a", "abbr", "acronym", "address", "area", "article", "aside", "audio", "b", "bdi", "bdo", "big", "blink", "blockquote", "body", "br", "button", "canvas", "caption", "center", "cite", "code", "col", "colgroup", "content", "data", "datalist", "dd", "decorator", "del", "details", "dfn", "dialog", "dir", "div", "dl", "dt", "element", "em", "fieldset", "figcaption", "figure", "font", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hgroup", "hr", "html", "i", "img", "input", "ins", "kbd", "label", "legend", "li", "main", "map", "mark", "marquee", "menu", "menuitem", "meter", "nav", "nobr", "ol", "optgroup", "option", "output", "p", "picture", "pre", "progress", "q", "rp", "rt", "ruby", "s", "samp", "section", "select", "shadow", "small", "source", "spacer", "span", "strike", "strong", "style", "sub", "summary", "sup", "table", "tbody", "td", "template", "textarea", "tfoot", "th", "thead", "time", "tr", "track", "tt", "u", "ul", "var", "video", "wbr"]);
  var svg = freeze(["svg", "a", "altglyph", "altglyphdef", "altglyphitem", "animatecolor", "animatemotion", "animatetransform", "circle", "clippath", "defs", "desc", "ellipse", "filter", "font", "g", "glyph", "glyphref", "hkern", "image", "line", "lineargradient", "marker", "mask", "metadata", "mpath", "path", "pattern", "polygon", "polyline", "radialgradient", "rect", "stop", "style", "switch", "symbol", "text", "textpath", "title", "tref", "tspan", "view", "vkern"]);
  var svgFilters = freeze(["feBlend", "feColorMatrix", "feComponentTransfer", "feComposite", "feConvolveMatrix", "feDiffuseLighting", "feDisplacementMap", "feDistantLight", "feFlood", "feFuncA", "feFuncB", "feFuncG", "feFuncR", "feGaussianBlur", "feMerge", "feMergeNode", "feMorphology", "feOffset", "fePointLight", "feSpecularLighting", "feSpotLight", "feTile", "feTurbulence"]);
  var svgDisallowed = freeze(["animate", "color-profile", "cursor", "discard", "fedropshadow", "feimage", "font-face", "font-face-format", "font-face-name", "font-face-src", "font-face-uri", "foreignobject", "hatch", "hatchpath", "mesh", "meshgradient", "meshpatch", "meshrow", "missing-glyph", "script", "set", "solidcolor", "unknown", "use"]);
  var mathMl = freeze(["math", "menclose", "merror", "mfenced", "mfrac", "mglyph", "mi", "mlabeledtr", "mmultiscripts", "mn", "mo", "mover", "mpadded", "mphantom", "mroot", "mrow", "ms", "mspace", "msqrt", "mstyle", "msub", "msup", "msubsup", "mtable", "mtd", "mtext", "mtr", "munder", "munderover"]);
  var mathMlDisallowed = freeze(["maction", "maligngroup", "malignmark", "mlongdiv", "mscarries", "mscarry", "msgroup", "mstack", "msline", "msrow", "semantics", "annotation", "annotation-xml", "mprescripts", "none"]);
  var text = freeze(["#text"]);
  var html$1 = freeze(["accept", "action", "align", "alt", "autocapitalize", "autocomplete", "autopictureinpicture", "autoplay", "background", "bgcolor", "border", "capture", "cellpadding", "cellspacing", "checked", "cite", "class", "clear", "color", "cols", "colspan", "controls", "controlslist", "coords", "crossorigin", "datetime", "decoding", "default", "dir", "disabled", "disablepictureinpicture", "disableremoteplayback", "download", "draggable", "enctype", "enterkeyhint", "face", "for", "headers", "height", "hidden", "high", "href", "hreflang", "id", "inputmode", "integrity", "ismap", "kind", "label", "lang", "list", "loading", "loop", "low", "max", "maxlength", "media", "method", "min", "minlength", "multiple", "muted", "name", "noshade", "novalidate", "nowrap", "open", "optimum", "pattern", "placeholder", "playsinline", "poster", "preload", "pubdate", "radiogroup", "readonly", "rel", "required", "rev", "reversed", "role", "rows", "rowspan", "spellcheck", "scope", "selected", "shape", "size", "sizes", "span", "srclang", "start", "src", "srcset", "step", "style", "summary", "tabindex", "title", "translate", "type", "usemap", "valign", "value", "width", "xmlns"]);
  var svg$1 = freeze(["accent-height", "accumulate", "additive", "alignment-baseline", "ascent", "attributename", "attributetype", "azimuth", "basefrequency", "baseline-shift", "begin", "bias", "by", "class", "clip", "clippathunits", "clip-path", "clip-rule", "color", "color-interpolation", "color-interpolation-filters", "color-profile", "color-rendering", "cx", "cy", "d", "dx", "dy", "diffuseconstant", "direction", "display", "divisor", "dur", "edgemode", "elevation", "end", "fill", "fill-opacity", "fill-rule", "filter", "filterunits", "flood-color", "flood-opacity", "font-family", "font-size", "font-size-adjust", "font-stretch", "font-style", "font-variant", "font-weight", "fx", "fy", "g1", "g2", "glyph-name", "glyphref", "gradientunits", "gradienttransform", "height", "href", "id", "image-rendering", "in", "in2", "k", "k1", "k2", "k3", "k4", "kerning", "keypoints", "keysplines", "keytimes", "lang", "lengthadjust", "letter-spacing", "kernelmatrix", "kernelunitlength", "lighting-color", "local", "marker-end", "marker-mid", "marker-start", "markerheight", "markerunits", "markerwidth", "maskcontentunits", "maskunits", "max", "mask", "media", "method", "mode", "min", "name", "numoctaves", "offset", "operator", "opacity", "order", "orient", "orientation", "origin", "overflow", "paint-order", "path", "pathlength", "patterncontentunits", "patterntransform", "patternunits", "points", "preservealpha", "preserveaspectratio", "primitiveunits", "r", "rx", "ry", "radius", "refx", "refy", "repeatcount", "repeatdur", "restart", "result", "rotate", "scale", "seed", "shape-rendering", "specularconstant", "specularexponent", "spreadmethod", "startoffset", "stddeviation", "stitchtiles", "stop-color", "stop-opacity", "stroke-dasharray", "stroke-dashoffset", "stroke-linecap", "stroke-linejoin", "stroke-miterlimit", "stroke-opacity", "stroke", "stroke-width", "style", "surfacescale", "systemlanguage", "tabindex", "targetx", "targety", "transform", "text-anchor", "text-decoration", "text-rendering", "textlength", "type", "u1", "u2", "unicode", "values", "viewbox", "visibility", "version", "vert-adv-y", "vert-origin-x", "vert-origin-y", "width", "word-spacing", "wrap", "writing-mode", "xchannelselector", "ychannelselector", "x", "x1", "x2", "xmlns", "y", "y1", "y2", "z", "zoomandpan"]);
  var mathMl$1 = freeze(["accent", "accentunder", "align", "bevelled", "close", "columnsalign", "columnlines", "columnspan", "denomalign", "depth", "dir", "display", "displaystyle", "encoding", "fence", "frame", "height", "href", "id", "largeop", "length", "linethickness", "lspace", "lquote", "mathbackground", "mathcolor", "mathsize", "mathvariant", "maxsize", "minsize", "movablelimits", "notation", "numalign", "open", "rowalign", "rowlines", "rowspacing", "rowspan", "rspace", "rquote", "scriptlevel", "scriptminsize", "scriptsizemultiplier", "selection", "separator", "separators", "stretchy", "subscriptshift", "supscriptshift", "symmetric", "voffset", "width", "xmlns"]);
  var xml = freeze(["xlink:href", "xml:id", "xlink:title", "xml:space", "xmlns:xlink"]);
  var MUSTACHE_EXPR = seal(/\{\{[\s\S]*|[\s\S]*\}\}/gm);
  var ERB_EXPR = seal(/<%[\s\S]*|[\s\S]*%>/gm);
  var DATA_ATTR = seal(/^data-[\-\w.\u00B7-\uFFFF]/);
  var ARIA_ATTR = seal(/^aria-[\-\w]+$/);
  var IS_ALLOWED_URI = seal(/^(?:(?:(?:f|ht)tps?|mailto|tel|callto|cid|xmpp):|[^a-z]|[a-z+.\-]+(?:[^a-z+.\-:]|$))/i);
  var IS_SCRIPT_OR_DATA = seal(/^(?:\w+script|data):/i);
  var ATTR_WHITESPACE = seal(/[\u0000-\u0020\u00A0\u1680\u180E\u2000-\u2029\u205F\u3000]/g);
  var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function(obj) {
    return typeof obj;
  } : function(obj) {
    return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;
  };
  function _toConsumableArray$1(arr) {
    if (Array.isArray(arr)) {
      for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) {
        arr2[i] = arr[i];
      }
      return arr2;
    } else {
      return Array.from(arr);
    }
  }
  var getGlobal = function getGlobal2() {
    return typeof window === "undefined" ? null : window;
  };
  var _createTrustedTypesPolicy = function _createTrustedTypesPolicy2(trustedTypes, document2) {
    if ((typeof trustedTypes === "undefined" ? "undefined" : _typeof(trustedTypes)) !== "object" || typeof trustedTypes.createPolicy !== "function") {
      return null;
    }
    var suffix = null;
    var ATTR_NAME = "data-tt-policy-suffix";
    if (document2.currentScript && document2.currentScript.hasAttribute(ATTR_NAME)) {
      suffix = document2.currentScript.getAttribute(ATTR_NAME);
    }
    var policyName = "dompurify" + (suffix ? "#" + suffix : "");
    try {
      return trustedTypes.createPolicy(policyName, {
        createHTML: function createHTML(html$$1) {
          return html$$1;
        }
      });
    } catch (_) {
      console.warn("TrustedTypes policy " + policyName + " could not be created.");
      return null;
    }
  };
  function createDOMPurify() {
    var window2 = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : getGlobal();
    var DOMPurify2 = function DOMPurify3(root) {
      return createDOMPurify(root);
    };
    DOMPurify2.version = "2.2.6";
    DOMPurify2.removed = [];
    if (!window2 || !window2.document || window2.document.nodeType !== 9) {
      DOMPurify2.isSupported = false;
      return DOMPurify2;
    }
    var originalDocument = window2.document;
    var document2 = window2.document;
    var DocumentFragment = window2.DocumentFragment, HTMLTemplateElement = window2.HTMLTemplateElement, Node = window2.Node, Element = window2.Element, NodeFilter = window2.NodeFilter, _window$NamedNodeMap = window2.NamedNodeMap, NamedNodeMap = _window$NamedNodeMap === void 0 ? window2.NamedNodeMap || window2.MozNamedAttrMap : _window$NamedNodeMap, Text = window2.Text, Comment = window2.Comment, DOMParser = window2.DOMParser, trustedTypes = window2.trustedTypes;
    var ElementPrototype = Element.prototype;
    var cloneNode = lookupGetter(ElementPrototype, "cloneNode");
    var getNextSibling = lookupGetter(ElementPrototype, "nextSibling");
    var getChildNodes = lookupGetter(ElementPrototype, "childNodes");
    var getParentNode = lookupGetter(ElementPrototype, "parentNode");
    if (typeof HTMLTemplateElement === "function") {
      var template = document2.createElement("template");
      if (template.content && template.content.ownerDocument) {
        document2 = template.content.ownerDocument;
      }
    }
    var trustedTypesPolicy = _createTrustedTypesPolicy(trustedTypes, originalDocument);
    var emptyHTML = trustedTypesPolicy && RETURN_TRUSTED_TYPE ? trustedTypesPolicy.createHTML("") : "";
    var _document = document2, implementation = _document.implementation, createNodeIterator = _document.createNodeIterator, getElementsByTagName = _document.getElementsByTagName, createDocumentFragment = _document.createDocumentFragment;
    var importNode = originalDocument.importNode;
    var documentMode = {};
    try {
      documentMode = clone(document2).documentMode ? document2.documentMode : {};
    } catch (_) {
    }
    var hooks = {};
    DOMPurify2.isSupported = implementation && typeof implementation.createHTMLDocument !== "undefined" && documentMode !== 9;
    var MUSTACHE_EXPR$$1 = MUSTACHE_EXPR, ERB_EXPR$$1 = ERB_EXPR, DATA_ATTR$$1 = DATA_ATTR, ARIA_ATTR$$1 = ARIA_ATTR, IS_SCRIPT_OR_DATA$$1 = IS_SCRIPT_OR_DATA, ATTR_WHITESPACE$$1 = ATTR_WHITESPACE;
    var IS_ALLOWED_URI$$1 = IS_ALLOWED_URI;
    var ALLOWED_TAGS = null;
    var DEFAULT_ALLOWED_TAGS = addToSet({}, [].concat(_toConsumableArray$1(html), _toConsumableArray$1(svg), _toConsumableArray$1(svgFilters), _toConsumableArray$1(mathMl), _toConsumableArray$1(text)));
    var ALLOWED_ATTR = null;
    var DEFAULT_ALLOWED_ATTR = addToSet({}, [].concat(_toConsumableArray$1(html$1), _toConsumableArray$1(svg$1), _toConsumableArray$1(mathMl$1), _toConsumableArray$1(xml)));
    var FORBID_TAGS = null;
    var FORBID_ATTR = null;
    var ALLOW_ARIA_ATTR = true;
    var ALLOW_DATA_ATTR = true;
    var ALLOW_UNKNOWN_PROTOCOLS = false;
    var SAFE_FOR_TEMPLATES = false;
    var WHOLE_DOCUMENT = false;
    var SET_CONFIG = false;
    var FORCE_BODY = false;
    var RETURN_DOM = false;
    var RETURN_DOM_FRAGMENT = false;
    var RETURN_DOM_IMPORT = true;
    var RETURN_TRUSTED_TYPE = false;
    var SANITIZE_DOM = true;
    var KEEP_CONTENT = true;
    var IN_PLACE = false;
    var USE_PROFILES = {};
    var FORBID_CONTENTS = addToSet({}, ["annotation-xml", "audio", "colgroup", "desc", "foreignobject", "head", "iframe", "math", "mi", "mn", "mo", "ms", "mtext", "noembed", "noframes", "noscript", "plaintext", "script", "style", "svg", "template", "thead", "title", "video", "xmp"]);
    var DATA_URI_TAGS = null;
    var DEFAULT_DATA_URI_TAGS = addToSet({}, ["audio", "video", "img", "source", "image", "track"]);
    var URI_SAFE_ATTRIBUTES = null;
    var DEFAULT_URI_SAFE_ATTRIBUTES = addToSet({}, ["alt", "class", "for", "id", "label", "name", "pattern", "placeholder", "summary", "title", "value", "style", "xmlns"]);
    var CONFIG = null;
    var formElement = document2.createElement("form");
    var _parseConfig = function _parseConfig2(cfg) {
      if (CONFIG && CONFIG === cfg) {
        return;
      }
      if (!cfg || (typeof cfg === "undefined" ? "undefined" : _typeof(cfg)) !== "object") {
        cfg = {};
      }
      cfg = clone(cfg);
      ALLOWED_TAGS = "ALLOWED_TAGS" in cfg ? addToSet({}, cfg.ALLOWED_TAGS) : DEFAULT_ALLOWED_TAGS;
      ALLOWED_ATTR = "ALLOWED_ATTR" in cfg ? addToSet({}, cfg.ALLOWED_ATTR) : DEFAULT_ALLOWED_ATTR;
      URI_SAFE_ATTRIBUTES = "ADD_URI_SAFE_ATTR" in cfg ? addToSet(clone(DEFAULT_URI_SAFE_ATTRIBUTES), cfg.ADD_URI_SAFE_ATTR) : DEFAULT_URI_SAFE_ATTRIBUTES;
      DATA_URI_TAGS = "ADD_DATA_URI_TAGS" in cfg ? addToSet(clone(DEFAULT_DATA_URI_TAGS), cfg.ADD_DATA_URI_TAGS) : DEFAULT_DATA_URI_TAGS;
      FORBID_TAGS = "FORBID_TAGS" in cfg ? addToSet({}, cfg.FORBID_TAGS) : {};
      FORBID_ATTR = "FORBID_ATTR" in cfg ? addToSet({}, cfg.FORBID_ATTR) : {};
      USE_PROFILES = "USE_PROFILES" in cfg ? cfg.USE_PROFILES : false;
      ALLOW_ARIA_ATTR = cfg.ALLOW_ARIA_ATTR !== false;
      ALLOW_DATA_ATTR = cfg.ALLOW_DATA_ATTR !== false;
      ALLOW_UNKNOWN_PROTOCOLS = cfg.ALLOW_UNKNOWN_PROTOCOLS || false;
      SAFE_FOR_TEMPLATES = cfg.SAFE_FOR_TEMPLATES || false;
      WHOLE_DOCUMENT = cfg.WHOLE_DOCUMENT || false;
      RETURN_DOM = cfg.RETURN_DOM || false;
      RETURN_DOM_FRAGMENT = cfg.RETURN_DOM_FRAGMENT || false;
      RETURN_DOM_IMPORT = cfg.RETURN_DOM_IMPORT !== false;
      RETURN_TRUSTED_TYPE = cfg.RETURN_TRUSTED_TYPE || false;
      FORCE_BODY = cfg.FORCE_BODY || false;
      SANITIZE_DOM = cfg.SANITIZE_DOM !== false;
      KEEP_CONTENT = cfg.KEEP_CONTENT !== false;
      IN_PLACE = cfg.IN_PLACE || false;
      IS_ALLOWED_URI$$1 = cfg.ALLOWED_URI_REGEXP || IS_ALLOWED_URI$$1;
      if (SAFE_FOR_TEMPLATES) {
        ALLOW_DATA_ATTR = false;
      }
      if (RETURN_DOM_FRAGMENT) {
        RETURN_DOM = true;
      }
      if (USE_PROFILES) {
        ALLOWED_TAGS = addToSet({}, [].concat(_toConsumableArray$1(text)));
        ALLOWED_ATTR = [];
        if (USE_PROFILES.html === true) {
          addToSet(ALLOWED_TAGS, html);
          addToSet(ALLOWED_ATTR, html$1);
        }
        if (USE_PROFILES.svg === true) {
          addToSet(ALLOWED_TAGS, svg);
          addToSet(ALLOWED_ATTR, svg$1);
          addToSet(ALLOWED_ATTR, xml);
        }
        if (USE_PROFILES.svgFilters === true) {
          addToSet(ALLOWED_TAGS, svgFilters);
          addToSet(ALLOWED_ATTR, svg$1);
          addToSet(ALLOWED_ATTR, xml);
        }
        if (USE_PROFILES.mathMl === true) {
          addToSet(ALLOWED_TAGS, mathMl);
          addToSet(ALLOWED_ATTR, mathMl$1);
          addToSet(ALLOWED_ATTR, xml);
        }
      }
      if (cfg.ADD_TAGS) {
        if (ALLOWED_TAGS === DEFAULT_ALLOWED_TAGS) {
          ALLOWED_TAGS = clone(ALLOWED_TAGS);
        }
        addToSet(ALLOWED_TAGS, cfg.ADD_TAGS);
      }
      if (cfg.ADD_ATTR) {
        if (ALLOWED_ATTR === DEFAULT_ALLOWED_ATTR) {
          ALLOWED_ATTR = clone(ALLOWED_ATTR);
        }
        addToSet(ALLOWED_ATTR, cfg.ADD_ATTR);
      }
      if (cfg.ADD_URI_SAFE_ATTR) {
        addToSet(URI_SAFE_ATTRIBUTES, cfg.ADD_URI_SAFE_ATTR);
      }
      if (KEEP_CONTENT) {
        ALLOWED_TAGS["#text"] = true;
      }
      if (WHOLE_DOCUMENT) {
        addToSet(ALLOWED_TAGS, ["html", "head", "body"]);
      }
      if (ALLOWED_TAGS.table) {
        addToSet(ALLOWED_TAGS, ["tbody"]);
        delete FORBID_TAGS.tbody;
      }
      if (freeze) {
        freeze(cfg);
      }
      CONFIG = cfg;
    };
    var MATHML_TEXT_INTEGRATION_POINTS = addToSet({}, ["mi", "mo", "mn", "ms", "mtext"]);
    var HTML_INTEGRATION_POINTS = addToSet({}, ["foreignobject", "desc", "title", "annotation-xml"]);
    var ALL_SVG_TAGS = addToSet({}, svg);
    addToSet(ALL_SVG_TAGS, svgFilters);
    addToSet(ALL_SVG_TAGS, svgDisallowed);
    var ALL_MATHML_TAGS = addToSet({}, mathMl);
    addToSet(ALL_MATHML_TAGS, mathMlDisallowed);
    var MATHML_NAMESPACE = "http://www.w3.org/1998/Math/MathML";
    var SVG_NAMESPACE = "http://www.w3.org/2000/svg";
    var HTML_NAMESPACE = "http://www.w3.org/1999/xhtml";
    var _checkValidNamespace = function _checkValidNamespace2(element) {
      var parent = getParentNode(element);
      if (!parent || !parent.tagName) {
        parent = {
          namespaceURI: HTML_NAMESPACE,
          tagName: "template"
        };
      }
      var tagName = stringToLowerCase(element.tagName);
      var parentTagName = stringToLowerCase(parent.tagName);
      if (element.namespaceURI === SVG_NAMESPACE) {
        if (parent.namespaceURI === HTML_NAMESPACE) {
          return tagName === "svg";
        }
        if (parent.namespaceURI === MATHML_NAMESPACE) {
          return tagName === "svg" && (parentTagName === "annotation-xml" || MATHML_TEXT_INTEGRATION_POINTS[parentTagName]);
        }
        return Boolean(ALL_SVG_TAGS[tagName]);
      }
      if (element.namespaceURI === MATHML_NAMESPACE) {
        if (parent.namespaceURI === HTML_NAMESPACE) {
          return tagName === "math";
        }
        if (parent.namespaceURI === SVG_NAMESPACE) {
          return tagName === "math" && HTML_INTEGRATION_POINTS[parentTagName];
        }
        return Boolean(ALL_MATHML_TAGS[tagName]);
      }
      if (element.namespaceURI === HTML_NAMESPACE) {
        if (parent.namespaceURI === SVG_NAMESPACE && !HTML_INTEGRATION_POINTS[parentTagName]) {
          return false;
        }
        if (parent.namespaceURI === MATHML_NAMESPACE && !MATHML_TEXT_INTEGRATION_POINTS[parentTagName]) {
          return false;
        }
        var commonSvgAndHTMLElements = addToSet({}, ["title", "style", "font", "a", "script"]);
        return !ALL_MATHML_TAGS[tagName] && (commonSvgAndHTMLElements[tagName] || !ALL_SVG_TAGS[tagName]);
      }
      return false;
    };
    var _forceRemove = function _forceRemove2(node) {
      arrayPush(DOMPurify2.removed, {element: node});
      try {
        node.parentNode.removeChild(node);
      } catch (_) {
        try {
          node.outerHTML = emptyHTML;
        } catch (_2) {
          node.remove();
        }
      }
    };
    var _removeAttribute = function _removeAttribute2(name, node) {
      try {
        arrayPush(DOMPurify2.removed, {
          attribute: node.getAttributeNode(name),
          from: node
        });
      } catch (_) {
        arrayPush(DOMPurify2.removed, {
          attribute: null,
          from: node
        });
      }
      node.removeAttribute(name);
    };
    var _initDocument = function _initDocument2(dirty) {
      var doc = void 0;
      var leadingWhitespace = void 0;
      if (FORCE_BODY) {
        dirty = "<remove></remove>" + dirty;
      } else {
        var matches = stringMatch(dirty, /^[\r\n\t ]+/);
        leadingWhitespace = matches && matches[0];
      }
      var dirtyPayload = trustedTypesPolicy ? trustedTypesPolicy.createHTML(dirty) : dirty;
      try {
        doc = new DOMParser().parseFromString(dirtyPayload, "text/html");
      } catch (_) {
      }
      if (!doc || !doc.documentElement) {
        doc = implementation.createHTMLDocument("");
        var _doc = doc, body = _doc.body;
        body.parentNode.removeChild(body.parentNode.firstElementChild);
        body.outerHTML = dirtyPayload;
      }
      if (dirty && leadingWhitespace) {
        doc.body.insertBefore(document2.createTextNode(leadingWhitespace), doc.body.childNodes[0] || null);
      }
      return getElementsByTagName.call(doc, WHOLE_DOCUMENT ? "html" : "body")[0];
    };
    var _createIterator = function _createIterator2(root) {
      return createNodeIterator.call(root.ownerDocument || root, root, NodeFilter.SHOW_ELEMENT | NodeFilter.SHOW_COMMENT | NodeFilter.SHOW_TEXT, function() {
        return NodeFilter.FILTER_ACCEPT;
      }, false);
    };
    var _isClobbered = function _isClobbered2(elm) {
      if (elm instanceof Text || elm instanceof Comment) {
        return false;
      }
      if (typeof elm.nodeName !== "string" || typeof elm.textContent !== "string" || typeof elm.removeChild !== "function" || !(elm.attributes instanceof NamedNodeMap) || typeof elm.removeAttribute !== "function" || typeof elm.setAttribute !== "function" || typeof elm.namespaceURI !== "string" || typeof elm.insertBefore !== "function") {
        return true;
      }
      return false;
    };
    var _isNode = function _isNode2(object) {
      return (typeof Node === "undefined" ? "undefined" : _typeof(Node)) === "object" ? object instanceof Node : object && (typeof object === "undefined" ? "undefined" : _typeof(object)) === "object" && typeof object.nodeType === "number" && typeof object.nodeName === "string";
    };
    var _executeHook = function _executeHook2(entryPoint, currentNode, data) {
      if (!hooks[entryPoint]) {
        return;
      }
      arrayForEach(hooks[entryPoint], function(hook) {
        hook.call(DOMPurify2, currentNode, data, CONFIG);
      });
    };
    var _sanitizeElements = function _sanitizeElements2(currentNode) {
      var content = void 0;
      _executeHook("beforeSanitizeElements", currentNode, null);
      if (_isClobbered(currentNode)) {
        _forceRemove(currentNode);
        return true;
      }
      if (stringMatch(currentNode.nodeName, /[\u0080-\uFFFF]/)) {
        _forceRemove(currentNode);
        return true;
      }
      var tagName = stringToLowerCase(currentNode.nodeName);
      _executeHook("uponSanitizeElement", currentNode, {
        tagName,
        allowedTags: ALLOWED_TAGS
      });
      if (!_isNode(currentNode.firstElementChild) && (!_isNode(currentNode.content) || !_isNode(currentNode.content.firstElementChild)) && regExpTest(/<[/\w]/g, currentNode.innerHTML) && regExpTest(/<[/\w]/g, currentNode.textContent)) {
        _forceRemove(currentNode);
        return true;
      }
      if (!ALLOWED_TAGS[tagName] || FORBID_TAGS[tagName]) {
        if (KEEP_CONTENT && !FORBID_CONTENTS[tagName]) {
          var parentNode = getParentNode(currentNode);
          var childNodes = getChildNodes(currentNode);
          var childCount = childNodes.length;
          for (var i = childCount - 1; i >= 0; --i) {
            parentNode.insertBefore(cloneNode(childNodes[i], true), getNextSibling(currentNode));
          }
        }
        _forceRemove(currentNode);
        return true;
      }
      if (currentNode instanceof Element && !_checkValidNamespace(currentNode)) {
        _forceRemove(currentNode);
        return true;
      }
      if ((tagName === "noscript" || tagName === "noembed") && regExpTest(/<\/no(script|embed)/i, currentNode.innerHTML)) {
        _forceRemove(currentNode);
        return true;
      }
      if (SAFE_FOR_TEMPLATES && currentNode.nodeType === 3) {
        content = currentNode.textContent;
        content = stringReplace(content, MUSTACHE_EXPR$$1, " ");
        content = stringReplace(content, ERB_EXPR$$1, " ");
        if (currentNode.textContent !== content) {
          arrayPush(DOMPurify2.removed, {element: currentNode.cloneNode()});
          currentNode.textContent = content;
        }
      }
      _executeHook("afterSanitizeElements", currentNode, null);
      return false;
    };
    var _isValidAttribute = function _isValidAttribute2(lcTag, lcName, value) {
      if (SANITIZE_DOM && (lcName === "id" || lcName === "name") && (value in document2 || value in formElement)) {
        return false;
      }
      if (ALLOW_DATA_ATTR && regExpTest(DATA_ATTR$$1, lcName))
        ;
      else if (ALLOW_ARIA_ATTR && regExpTest(ARIA_ATTR$$1, lcName))
        ;
      else if (!ALLOWED_ATTR[lcName] || FORBID_ATTR[lcName]) {
        return false;
      } else if (URI_SAFE_ATTRIBUTES[lcName])
        ;
      else if (regExpTest(IS_ALLOWED_URI$$1, stringReplace(value, ATTR_WHITESPACE$$1, "")))
        ;
      else if ((lcName === "src" || lcName === "xlink:href" || lcName === "href") && lcTag !== "script" && stringIndexOf(value, "data:") === 0 && DATA_URI_TAGS[lcTag])
        ;
      else if (ALLOW_UNKNOWN_PROTOCOLS && !regExpTest(IS_SCRIPT_OR_DATA$$1, stringReplace(value, ATTR_WHITESPACE$$1, "")))
        ;
      else if (!value)
        ;
      else {
        return false;
      }
      return true;
    };
    var _sanitizeAttributes = function _sanitizeAttributes2(currentNode) {
      var attr = void 0;
      var value = void 0;
      var lcName = void 0;
      var l = void 0;
      _executeHook("beforeSanitizeAttributes", currentNode, null);
      var attributes = currentNode.attributes;
      if (!attributes) {
        return;
      }
      var hookEvent = {
        attrName: "",
        attrValue: "",
        keepAttr: true,
        allowedAttributes: ALLOWED_ATTR
      };
      l = attributes.length;
      while (l--) {
        attr = attributes[l];
        var _attr = attr, name = _attr.name, namespaceURI = _attr.namespaceURI;
        value = stringTrim(attr.value);
        lcName = stringToLowerCase(name);
        hookEvent.attrName = lcName;
        hookEvent.attrValue = value;
        hookEvent.keepAttr = true;
        hookEvent.forceKeepAttr = void 0;
        _executeHook("uponSanitizeAttribute", currentNode, hookEvent);
        value = hookEvent.attrValue;
        if (hookEvent.forceKeepAttr) {
          continue;
        }
        _removeAttribute(name, currentNode);
        if (!hookEvent.keepAttr) {
          continue;
        }
        if (regExpTest(/\/>/i, value)) {
          _removeAttribute(name, currentNode);
          continue;
        }
        if (SAFE_FOR_TEMPLATES) {
          value = stringReplace(value, MUSTACHE_EXPR$$1, " ");
          value = stringReplace(value, ERB_EXPR$$1, " ");
        }
        var lcTag = currentNode.nodeName.toLowerCase();
        if (!_isValidAttribute(lcTag, lcName, value)) {
          continue;
        }
        try {
          if (namespaceURI) {
            currentNode.setAttributeNS(namespaceURI, name, value);
          } else {
            currentNode.setAttribute(name, value);
          }
          arrayPop(DOMPurify2.removed);
        } catch (_) {
        }
      }
      _executeHook("afterSanitizeAttributes", currentNode, null);
    };
    var _sanitizeShadowDOM = function _sanitizeShadowDOM2(fragment) {
      var shadowNode = void 0;
      var shadowIterator = _createIterator(fragment);
      _executeHook("beforeSanitizeShadowDOM", fragment, null);
      while (shadowNode = shadowIterator.nextNode()) {
        _executeHook("uponSanitizeShadowNode", shadowNode, null);
        if (_sanitizeElements(shadowNode)) {
          continue;
        }
        if (shadowNode.content instanceof DocumentFragment) {
          _sanitizeShadowDOM2(shadowNode.content);
        }
        _sanitizeAttributes(shadowNode);
      }
      _executeHook("afterSanitizeShadowDOM", fragment, null);
    };
    DOMPurify2.sanitize = function(dirty, cfg) {
      var body = void 0;
      var importedNode = void 0;
      var currentNode = void 0;
      var oldNode = void 0;
      var returnNode = void 0;
      if (!dirty) {
        dirty = "<!-->";
      }
      if (typeof dirty !== "string" && !_isNode(dirty)) {
        if (typeof dirty.toString !== "function") {
          throw typeErrorCreate("toString is not a function");
        } else {
          dirty = dirty.toString();
          if (typeof dirty !== "string") {
            throw typeErrorCreate("dirty is not a string, aborting");
          }
        }
      }
      if (!DOMPurify2.isSupported) {
        if (_typeof(window2.toStaticHTML) === "object" || typeof window2.toStaticHTML === "function") {
          if (typeof dirty === "string") {
            return window2.toStaticHTML(dirty);
          }
          if (_isNode(dirty)) {
            return window2.toStaticHTML(dirty.outerHTML);
          }
        }
        return dirty;
      }
      if (!SET_CONFIG) {
        _parseConfig(cfg);
      }
      DOMPurify2.removed = [];
      if (typeof dirty === "string") {
        IN_PLACE = false;
      }
      if (IN_PLACE)
        ;
      else if (dirty instanceof Node) {
        body = _initDocument("<!---->");
        importedNode = body.ownerDocument.importNode(dirty, true);
        if (importedNode.nodeType === 1 && importedNode.nodeName === "BODY") {
          body = importedNode;
        } else if (importedNode.nodeName === "HTML") {
          body = importedNode;
        } else {
          body.appendChild(importedNode);
        }
      } else {
        if (!RETURN_DOM && !SAFE_FOR_TEMPLATES && !WHOLE_DOCUMENT && dirty.indexOf("<") === -1) {
          return trustedTypesPolicy && RETURN_TRUSTED_TYPE ? trustedTypesPolicy.createHTML(dirty) : dirty;
        }
        body = _initDocument(dirty);
        if (!body) {
          return RETURN_DOM ? null : emptyHTML;
        }
      }
      if (body && FORCE_BODY) {
        _forceRemove(body.firstChild);
      }
      var nodeIterator = _createIterator(IN_PLACE ? dirty : body);
      while (currentNode = nodeIterator.nextNode()) {
        if (currentNode.nodeType === 3 && currentNode === oldNode) {
          continue;
        }
        if (_sanitizeElements(currentNode)) {
          continue;
        }
        if (currentNode.content instanceof DocumentFragment) {
          _sanitizeShadowDOM(currentNode.content);
        }
        _sanitizeAttributes(currentNode);
        oldNode = currentNode;
      }
      oldNode = null;
      if (IN_PLACE) {
        return dirty;
      }
      if (RETURN_DOM) {
        if (RETURN_DOM_FRAGMENT) {
          returnNode = createDocumentFragment.call(body.ownerDocument);
          while (body.firstChild) {
            returnNode.appendChild(body.firstChild);
          }
        } else {
          returnNode = body;
        }
        if (RETURN_DOM_IMPORT) {
          returnNode = importNode.call(originalDocument, returnNode, true);
        }
        return returnNode;
      }
      var serializedHTML = WHOLE_DOCUMENT ? body.outerHTML : body.innerHTML;
      if (SAFE_FOR_TEMPLATES) {
        serializedHTML = stringReplace(serializedHTML, MUSTACHE_EXPR$$1, " ");
        serializedHTML = stringReplace(serializedHTML, ERB_EXPR$$1, " ");
      }
      return trustedTypesPolicy && RETURN_TRUSTED_TYPE ? trustedTypesPolicy.createHTML(serializedHTML) : serializedHTML;
    };
    DOMPurify2.setConfig = function(cfg) {
      _parseConfig(cfg);
      SET_CONFIG = true;
    };
    DOMPurify2.clearConfig = function() {
      CONFIG = null;
      SET_CONFIG = false;
    };
    DOMPurify2.isValidAttribute = function(tag, attr, value) {
      if (!CONFIG) {
        _parseConfig({});
      }
      var lcTag = stringToLowerCase(tag);
      var lcName = stringToLowerCase(attr);
      return _isValidAttribute(lcTag, lcName, value);
    };
    DOMPurify2.addHook = function(entryPoint, hookFunction) {
      if (typeof hookFunction !== "function") {
        return;
      }
      hooks[entryPoint] = hooks[entryPoint] || [];
      arrayPush(hooks[entryPoint], hookFunction);
    };
    DOMPurify2.removeHook = function(entryPoint) {
      if (hooks[entryPoint]) {
        arrayPop(hooks[entryPoint]);
      }
    };
    DOMPurify2.removeHooks = function(entryPoint) {
      if (hooks[entryPoint]) {
        hooks[entryPoint] = [];
      }
    };
    DOMPurify2.removeAllHooks = function() {
      hooks = {};
    };
    return DOMPurify2;
  }
  var purify = createDOMPurify();
  module2.exports = purify;
});

// node_modules/scheduler/cjs/scheduler.production.min.js
var require_scheduler_production_min = __commonJS((exports2) => {
  /** @license React v0.20.1
   * scheduler.production.min.js
   *
   * Copyright (c) Facebook, Inc. and its affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   */
  "use strict";
  var f;
  var g;
  var h;
  var k;
  if (typeof performance === "object" && typeof performance.now === "function") {
    l = performance;
    exports2.unstable_now = function() {
      return l.now();
    };
  } else {
    p = Date, q = p.now();
    exports2.unstable_now = function() {
      return p.now() - q;
    };
  }
  var l;
  var p;
  var q;
  if (typeof window === "undefined" || typeof MessageChannel !== "function") {
    t = null, u = null, w = function() {
      if (t !== null)
        try {
          var a = exports2.unstable_now();
          t(true, a);
          t = null;
        } catch (b) {
          throw setTimeout(w, 0), b;
        }
    };
    f = function(a) {
      t !== null ? setTimeout(f, 0, a) : (t = a, setTimeout(w, 0));
    };
    g = function(a, b) {
      u = setTimeout(a, b);
    };
    h = function() {
      clearTimeout(u);
    };
    exports2.unstable_shouldYield = function() {
      return false;
    };
    k = exports2.unstable_forceFrameRate = function() {
    };
  } else {
    x = window.setTimeout, y = window.clearTimeout;
    if (typeof console !== "undefined") {
      z = window.cancelAnimationFrame;
      typeof window.requestAnimationFrame !== "function" && console.error("This browser doesn't support requestAnimationFrame. Make sure that you load a polyfill in older browsers. https://reactjs.org/link/react-polyfills");
      typeof z !== "function" && console.error("This browser doesn't support cancelAnimationFrame. Make sure that you load a polyfill in older browsers. https://reactjs.org/link/react-polyfills");
    }
    A = false, B = null, C = -1, D = 5, E = 0;
    exports2.unstable_shouldYield = function() {
      return exports2.unstable_now() >= E;
    };
    k = function() {
    };
    exports2.unstable_forceFrameRate = function(a) {
      0 > a || 125 < a ? console.error("forceFrameRate takes a positive int between 0 and 125, forcing frame rates higher than 125 fps is not supported") : D = 0 < a ? Math.floor(1e3 / a) : 5;
    };
    F = new MessageChannel(), G = F.port2;
    F.port1.onmessage = function() {
      if (B !== null) {
        var a = exports2.unstable_now();
        E = a + D;
        try {
          B(true, a) ? G.postMessage(null) : (A = false, B = null);
        } catch (b) {
          throw G.postMessage(null), b;
        }
      } else
        A = false;
    };
    f = function(a) {
      B = a;
      A || (A = true, G.postMessage(null));
    };
    g = function(a, b) {
      C = x(function() {
        a(exports2.unstable_now());
      }, b);
    };
    h = function() {
      y(C);
      C = -1;
    };
  }
  var t;
  var u;
  var w;
  var x;
  var y;
  var z;
  var A;
  var B;
  var C;
  var D;
  var E;
  var F;
  var G;
  function H(a, b) {
    var c = a.length;
    a.push(b);
    a:
      for (; ; ) {
        var d = c - 1 >>> 1, e = a[d];
        if (e !== void 0 && 0 < I(e, b))
          a[d] = b, a[c] = e, c = d;
        else
          break a;
      }
  }
  function J(a) {
    a = a[0];
    return a === void 0 ? null : a;
  }
  function K(a) {
    var b = a[0];
    if (b !== void 0) {
      var c = a.pop();
      if (c !== b) {
        a[0] = c;
        a:
          for (var d = 0, e = a.length; d < e; ) {
            var m = 2 * (d + 1) - 1, n = a[m], v = m + 1, r = a[v];
            if (n !== void 0 && 0 > I(n, c))
              r !== void 0 && 0 > I(r, n) ? (a[d] = r, a[v] = c, d = v) : (a[d] = n, a[m] = c, d = m);
            else if (r !== void 0 && 0 > I(r, c))
              a[d] = r, a[v] = c, d = v;
            else
              break a;
          }
      }
      return b;
    }
    return null;
  }
  function I(a, b) {
    var c = a.sortIndex - b.sortIndex;
    return c !== 0 ? c : a.id - b.id;
  }
  var L = [];
  var M = [];
  var N = 1;
  var O = null;
  var P = 3;
  var Q = false;
  var R = false;
  var S = false;
  function T(a) {
    for (var b = J(M); b !== null; ) {
      if (b.callback === null)
        K(M);
      else if (b.startTime <= a)
        K(M), b.sortIndex = b.expirationTime, H(L, b);
      else
        break;
      b = J(M);
    }
  }
  function U(a) {
    S = false;
    T(a);
    if (!R)
      if (J(L) !== null)
        R = true, f(V);
      else {
        var b = J(M);
        b !== null && g(U, b.startTime - a);
      }
  }
  function V(a, b) {
    R = false;
    S && (S = false, h());
    Q = true;
    var c = P;
    try {
      T(b);
      for (O = J(L); O !== null && (!(O.expirationTime > b) || a && !exports2.unstable_shouldYield()); ) {
        var d = O.callback;
        if (typeof d === "function") {
          O.callback = null;
          P = O.priorityLevel;
          var e = d(O.expirationTime <= b);
          b = exports2.unstable_now();
          typeof e === "function" ? O.callback = e : O === J(L) && K(L);
          T(b);
        } else
          K(L);
        O = J(L);
      }
      if (O !== null)
        var m = true;
      else {
        var n = J(M);
        n !== null && g(U, n.startTime - b);
        m = false;
      }
      return m;
    } finally {
      O = null, P = c, Q = false;
    }
  }
  var W = k;
  exports2.unstable_IdlePriority = 5;
  exports2.unstable_ImmediatePriority = 1;
  exports2.unstable_LowPriority = 4;
  exports2.unstable_NormalPriority = 3;
  exports2.unstable_Profiling = null;
  exports2.unstable_UserBlockingPriority = 2;
  exports2.unstable_cancelCallback = function(a) {
    a.callback = null;
  };
  exports2.unstable_continueExecution = function() {
    R || Q || (R = true, f(V));
  };
  exports2.unstable_getCurrentPriorityLevel = function() {
    return P;
  };
  exports2.unstable_getFirstCallbackNode = function() {
    return J(L);
  };
  exports2.unstable_next = function(a) {
    switch (P) {
      case 1:
      case 2:
      case 3:
        var b = 3;
        break;
      default:
        b = P;
    }
    var c = P;
    P = b;
    try {
      return a();
    } finally {
      P = c;
    }
  };
  exports2.unstable_pauseExecution = function() {
  };
  exports2.unstable_requestPaint = W;
  exports2.unstable_runWithPriority = function(a, b) {
    switch (a) {
      case 1:
      case 2:
      case 3:
      case 4:
      case 5:
        break;
      default:
        a = 3;
    }
    var c = P;
    P = a;
    try {
      return b();
    } finally {
      P = c;
    }
  };
  exports2.unstable_scheduleCallback = function(a, b, c) {
    var d = exports2.unstable_now();
    typeof c === "object" && c !== null ? (c = c.delay, c = typeof c === "number" && 0 < c ? d + c : d) : c = d;
    switch (a) {
      case 1:
        var e = -1;
        break;
      case 2:
        e = 250;
        break;
      case 5:
        e = 1073741823;
        break;
      case 4:
        e = 1e4;
        break;
      default:
        e = 5e3;
    }
    e = c + e;
    a = {id: N++, callback: b, priorityLevel: a, startTime: c, expirationTime: e, sortIndex: -1};
    c > d ? (a.sortIndex = c, H(M, a), J(L) === null && a === J(M) && (S ? h() : S = true, g(U, c - d))) : (a.sortIndex = e, H(L, a), R || Q || (R = true, f(V)));
    return a;
  };
  exports2.unstable_wrapCallback = function(a) {
    var b = P;
    return function() {
      var c = P;
      P = b;
      try {
        return a.apply(this, arguments);
      } finally {
        P = c;
      }
    };
  };
});

// node_modules/scheduler/index.js
var require_scheduler = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_scheduler_production_min();
  } else {
    module2.exports = null;
  }
});

// node_modules/react-dom/cjs/react-dom.production.min.js
var require_react_dom_production_min = __commonJS((exports2) => {
  /** @license React v17.0.1
   * react-dom.production.min.js
   *
   * Copyright (c) Facebook, Inc. and its affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   */
  "use strict";
  var aa = require_react();
  var m = require_object_assign();
  var r = require_scheduler();
  function y(a) {
    for (var b = "https://reactjs.org/docs/error-decoder.html?invariant=" + a, c = 1; c < arguments.length; c++)
      b += "&args[]=" + encodeURIComponent(arguments[c]);
    return "Minified React error #" + a + "; visit " + b + " for the full message or use the non-minified dev environment for full errors and additional helpful warnings.";
  }
  if (!aa)
    throw Error(y(227));
  var ba = new Set();
  var ca = {};
  function da(a, b) {
    ea(a, b);
    ea(a + "Capture", b);
  }
  function ea(a, b) {
    ca[a] = b;
    for (a = 0; a < b.length; a++)
      ba.add(b[a]);
  }
  var fa = !(typeof window === "undefined" || typeof window.document === "undefined" || typeof window.document.createElement === "undefined");
  var ha = /^[:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD][:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD\-.0-9\u00B7\u0300-\u036F\u203F-\u2040]*$/;
  var ia = Object.prototype.hasOwnProperty;
  var ja = {};
  var ka = {};
  function la(a) {
    if (ia.call(ka, a))
      return true;
    if (ia.call(ja, a))
      return false;
    if (ha.test(a))
      return ka[a] = true;
    ja[a] = true;
    return false;
  }
  function ma(a, b, c, d) {
    if (c !== null && c.type === 0)
      return false;
    switch (typeof b) {
      case "function":
      case "symbol":
        return true;
      case "boolean":
        if (d)
          return false;
        if (c !== null)
          return !c.acceptsBooleans;
        a = a.toLowerCase().slice(0, 5);
        return a !== "data-" && a !== "aria-";
      default:
        return false;
    }
  }
  function na(a, b, c, d) {
    if (b === null || typeof b === "undefined" || ma(a, b, c, d))
      return true;
    if (d)
      return false;
    if (c !== null)
      switch (c.type) {
        case 3:
          return !b;
        case 4:
          return b === false;
        case 5:
          return isNaN(b);
        case 6:
          return isNaN(b) || 1 > b;
      }
    return false;
  }
  function B(a, b, c, d, e, f, g) {
    this.acceptsBooleans = b === 2 || b === 3 || b === 4;
    this.attributeName = d;
    this.attributeNamespace = e;
    this.mustUseProperty = c;
    this.propertyName = a;
    this.type = b;
    this.sanitizeURL = f;
    this.removeEmptyString = g;
  }
  var D = {};
  "children dangerouslySetInnerHTML defaultValue defaultChecked innerHTML suppressContentEditableWarning suppressHydrationWarning style".split(" ").forEach(function(a) {
    D[a] = new B(a, 0, false, a, null, false, false);
  });
  [["acceptCharset", "accept-charset"], ["className", "class"], ["htmlFor", "for"], ["httpEquiv", "http-equiv"]].forEach(function(a) {
    var b = a[0];
    D[b] = new B(b, 1, false, a[1], null, false, false);
  });
  ["contentEditable", "draggable", "spellCheck", "value"].forEach(function(a) {
    D[a] = new B(a, 2, false, a.toLowerCase(), null, false, false);
  });
  ["autoReverse", "externalResourcesRequired", "focusable", "preserveAlpha"].forEach(function(a) {
    D[a] = new B(a, 2, false, a, null, false, false);
  });
  "allowFullScreen async autoFocus autoPlay controls default defer disabled disablePictureInPicture disableRemotePlayback formNoValidate hidden loop noModule noValidate open playsInline readOnly required reversed scoped seamless itemScope".split(" ").forEach(function(a) {
    D[a] = new B(a, 3, false, a.toLowerCase(), null, false, false);
  });
  ["checked", "multiple", "muted", "selected"].forEach(function(a) {
    D[a] = new B(a, 3, true, a, null, false, false);
  });
  ["capture", "download"].forEach(function(a) {
    D[a] = new B(a, 4, false, a, null, false, false);
  });
  ["cols", "rows", "size", "span"].forEach(function(a) {
    D[a] = new B(a, 6, false, a, null, false, false);
  });
  ["rowSpan", "start"].forEach(function(a) {
    D[a] = new B(a, 5, false, a.toLowerCase(), null, false, false);
  });
  var oa = /[\-:]([a-z])/g;
  function pa(a) {
    return a[1].toUpperCase();
  }
  "accent-height alignment-baseline arabic-form baseline-shift cap-height clip-path clip-rule color-interpolation color-interpolation-filters color-profile color-rendering dominant-baseline enable-background fill-opacity fill-rule flood-color flood-opacity font-family font-size font-size-adjust font-stretch font-style font-variant font-weight glyph-name glyph-orientation-horizontal glyph-orientation-vertical horiz-adv-x horiz-origin-x image-rendering letter-spacing lighting-color marker-end marker-mid marker-start overline-position overline-thickness paint-order panose-1 pointer-events rendering-intent shape-rendering stop-color stop-opacity strikethrough-position strikethrough-thickness stroke-dasharray stroke-dashoffset stroke-linecap stroke-linejoin stroke-miterlimit stroke-opacity stroke-width text-anchor text-decoration text-rendering underline-position underline-thickness unicode-bidi unicode-range units-per-em v-alphabetic v-hanging v-ideographic v-mathematical vector-effect vert-adv-y vert-origin-x vert-origin-y word-spacing writing-mode xmlns:xlink x-height".split(" ").forEach(function(a) {
    var b = a.replace(oa, pa);
    D[b] = new B(b, 1, false, a, null, false, false);
  });
  "xlink:actuate xlink:arcrole xlink:role xlink:show xlink:title xlink:type".split(" ").forEach(function(a) {
    var b = a.replace(oa, pa);
    D[b] = new B(b, 1, false, a, "http://www.w3.org/1999/xlink", false, false);
  });
  ["xml:base", "xml:lang", "xml:space"].forEach(function(a) {
    var b = a.replace(oa, pa);
    D[b] = new B(b, 1, false, a, "http://www.w3.org/XML/1998/namespace", false, false);
  });
  ["tabIndex", "crossOrigin"].forEach(function(a) {
    D[a] = new B(a, 1, false, a.toLowerCase(), null, false, false);
  });
  D.xlinkHref = new B("xlinkHref", 1, false, "xlink:href", "http://www.w3.org/1999/xlink", true, false);
  ["src", "href", "action", "formAction"].forEach(function(a) {
    D[a] = new B(a, 1, false, a.toLowerCase(), null, true, true);
  });
  function qa(a, b, c, d) {
    var e = D.hasOwnProperty(b) ? D[b] : null;
    var f = e !== null ? e.type === 0 : d ? false : !(2 < b.length) || b[0] !== "o" && b[0] !== "O" || b[1] !== "n" && b[1] !== "N" ? false : true;
    f || (na(b, c, e, d) && (c = null), d || e === null ? la(b) && (c === null ? a.removeAttribute(b) : a.setAttribute(b, "" + c)) : e.mustUseProperty ? a[e.propertyName] = c === null ? e.type === 3 ? false : "" : c : (b = e.attributeName, d = e.attributeNamespace, c === null ? a.removeAttribute(b) : (e = e.type, c = e === 3 || e === 4 && c === true ? "" : "" + c, d ? a.setAttributeNS(d, b, c) : a.setAttribute(b, c))));
  }
  var ra = aa.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED;
  var sa = 60103;
  var ta = 60106;
  var ua = 60107;
  var wa = 60108;
  var xa = 60114;
  var ya = 60109;
  var za = 60110;
  var Aa = 60112;
  var Ba = 60113;
  var Ca = 60120;
  var Da = 60115;
  var Ea = 60116;
  var Fa = 60121;
  var Ga = 60128;
  var Ha = 60129;
  var Ia = 60130;
  var Ja = 60131;
  if (typeof Symbol === "function" && Symbol.for) {
    E = Symbol.for;
    sa = E("react.element");
    ta = E("react.portal");
    ua = E("react.fragment");
    wa = E("react.strict_mode");
    xa = E("react.profiler");
    ya = E("react.provider");
    za = E("react.context");
    Aa = E("react.forward_ref");
    Ba = E("react.suspense");
    Ca = E("react.suspense_list");
    Da = E("react.memo");
    Ea = E("react.lazy");
    Fa = E("react.block");
    E("react.scope");
    Ga = E("react.opaque.id");
    Ha = E("react.debug_trace_mode");
    Ia = E("react.offscreen");
    Ja = E("react.legacy_hidden");
  }
  var E;
  var Ka = typeof Symbol === "function" && Symbol.iterator;
  function La(a) {
    if (a === null || typeof a !== "object")
      return null;
    a = Ka && a[Ka] || a["@@iterator"];
    return typeof a === "function" ? a : null;
  }
  var Ma;
  function Na(a) {
    if (Ma === void 0)
      try {
        throw Error();
      } catch (c) {
        var b = c.stack.trim().match(/\n( *(at )?)/);
        Ma = b && b[1] || "";
      }
    return "\n" + Ma + a;
  }
  var Oa = false;
  function Pa(a, b) {
    if (!a || Oa)
      return "";
    Oa = true;
    var c = Error.prepareStackTrace;
    Error.prepareStackTrace = void 0;
    try {
      if (b)
        if (b = function() {
          throw Error();
        }, Object.defineProperty(b.prototype, "props", {set: function() {
          throw Error();
        }}), typeof Reflect === "object" && Reflect.construct) {
          try {
            Reflect.construct(b, []);
          } catch (k) {
            var d = k;
          }
          Reflect.construct(a, [], b);
        } else {
          try {
            b.call();
          } catch (k) {
            d = k;
          }
          a.call(b.prototype);
        }
      else {
        try {
          throw Error();
        } catch (k) {
          d = k;
        }
        a();
      }
    } catch (k) {
      if (k && d && typeof k.stack === "string") {
        for (var e = k.stack.split("\n"), f = d.stack.split("\n"), g = e.length - 1, h = f.length - 1; 1 <= g && 0 <= h && e[g] !== f[h]; )
          h--;
        for (; 1 <= g && 0 <= h; g--, h--)
          if (e[g] !== f[h]) {
            if (g !== 1 || h !== 1) {
              do
                if (g--, h--, 0 > h || e[g] !== f[h])
                  return "\n" + e[g].replace(" at new ", " at ");
              while (1 <= g && 0 <= h);
            }
            break;
          }
      }
    } finally {
      Oa = false, Error.prepareStackTrace = c;
    }
    return (a = a ? a.displayName || a.name : "") ? Na(a) : "";
  }
  function Qa(a) {
    switch (a.tag) {
      case 5:
        return Na(a.type);
      case 16:
        return Na("Lazy");
      case 13:
        return Na("Suspense");
      case 19:
        return Na("SuspenseList");
      case 0:
      case 2:
      case 15:
        return a = Pa(a.type, false), a;
      case 11:
        return a = Pa(a.type.render, false), a;
      case 22:
        return a = Pa(a.type._render, false), a;
      case 1:
        return a = Pa(a.type, true), a;
      default:
        return "";
    }
  }
  function Ra(a) {
    if (a == null)
      return null;
    if (typeof a === "function")
      return a.displayName || a.name || null;
    if (typeof a === "string")
      return a;
    switch (a) {
      case ua:
        return "Fragment";
      case ta:
        return "Portal";
      case xa:
        return "Profiler";
      case wa:
        return "StrictMode";
      case Ba:
        return "Suspense";
      case Ca:
        return "SuspenseList";
    }
    if (typeof a === "object")
      switch (a.$$typeof) {
        case za:
          return (a.displayName || "Context") + ".Consumer";
        case ya:
          return (a._context.displayName || "Context") + ".Provider";
        case Aa:
          var b = a.render;
          b = b.displayName || b.name || "";
          return a.displayName || (b !== "" ? "ForwardRef(" + b + ")" : "ForwardRef");
        case Da:
          return Ra(a.type);
        case Fa:
          return Ra(a._render);
        case Ea:
          b = a._payload;
          a = a._init;
          try {
            return Ra(a(b));
          } catch (c) {
          }
      }
    return null;
  }
  function Sa(a) {
    switch (typeof a) {
      case "boolean":
      case "number":
      case "object":
      case "string":
      case "undefined":
        return a;
      default:
        return "";
    }
  }
  function Ta(a) {
    var b = a.type;
    return (a = a.nodeName) && a.toLowerCase() === "input" && (b === "checkbox" || b === "radio");
  }
  function Ua(a) {
    var b = Ta(a) ? "checked" : "value", c = Object.getOwnPropertyDescriptor(a.constructor.prototype, b), d = "" + a[b];
    if (!a.hasOwnProperty(b) && typeof c !== "undefined" && typeof c.get === "function" && typeof c.set === "function") {
      var e = c.get, f = c.set;
      Object.defineProperty(a, b, {configurable: true, get: function() {
        return e.call(this);
      }, set: function(a2) {
        d = "" + a2;
        f.call(this, a2);
      }});
      Object.defineProperty(a, b, {enumerable: c.enumerable});
      return {getValue: function() {
        return d;
      }, setValue: function(a2) {
        d = "" + a2;
      }, stopTracking: function() {
        a._valueTracker = null;
        delete a[b];
      }};
    }
  }
  function Va(a) {
    a._valueTracker || (a._valueTracker = Ua(a));
  }
  function Wa(a) {
    if (!a)
      return false;
    var b = a._valueTracker;
    if (!b)
      return true;
    var c = b.getValue();
    var d = "";
    a && (d = Ta(a) ? a.checked ? "true" : "false" : a.value);
    a = d;
    return a !== c ? (b.setValue(a), true) : false;
  }
  function Xa(a) {
    a = a || (typeof document !== "undefined" ? document : void 0);
    if (typeof a === "undefined")
      return null;
    try {
      return a.activeElement || a.body;
    } catch (b) {
      return a.body;
    }
  }
  function Ya(a, b) {
    var c = b.checked;
    return m({}, b, {defaultChecked: void 0, defaultValue: void 0, value: void 0, checked: c != null ? c : a._wrapperState.initialChecked});
  }
  function Za(a, b) {
    var c = b.defaultValue == null ? "" : b.defaultValue, d = b.checked != null ? b.checked : b.defaultChecked;
    c = Sa(b.value != null ? b.value : c);
    a._wrapperState = {initialChecked: d, initialValue: c, controlled: b.type === "checkbox" || b.type === "radio" ? b.checked != null : b.value != null};
  }
  function $a(a, b) {
    b = b.checked;
    b != null && qa(a, "checked", b, false);
  }
  function ab(a, b) {
    $a(a, b);
    var c = Sa(b.value), d = b.type;
    if (c != null)
      if (d === "number") {
        if (c === 0 && a.value === "" || a.value != c)
          a.value = "" + c;
      } else
        a.value !== "" + c && (a.value = "" + c);
    else if (d === "submit" || d === "reset") {
      a.removeAttribute("value");
      return;
    }
    b.hasOwnProperty("value") ? bb(a, b.type, c) : b.hasOwnProperty("defaultValue") && bb(a, b.type, Sa(b.defaultValue));
    b.checked == null && b.defaultChecked != null && (a.defaultChecked = !!b.defaultChecked);
  }
  function cb(a, b, c) {
    if (b.hasOwnProperty("value") || b.hasOwnProperty("defaultValue")) {
      var d = b.type;
      if (!(d !== "submit" && d !== "reset" || b.value !== void 0 && b.value !== null))
        return;
      b = "" + a._wrapperState.initialValue;
      c || b === a.value || (a.value = b);
      a.defaultValue = b;
    }
    c = a.name;
    c !== "" && (a.name = "");
    a.defaultChecked = !!a._wrapperState.initialChecked;
    c !== "" && (a.name = c);
  }
  function bb(a, b, c) {
    if (b !== "number" || Xa(a.ownerDocument) !== a)
      c == null ? a.defaultValue = "" + a._wrapperState.initialValue : a.defaultValue !== "" + c && (a.defaultValue = "" + c);
  }
  function db(a) {
    var b = "";
    aa.Children.forEach(a, function(a2) {
      a2 != null && (b += a2);
    });
    return b;
  }
  function eb(a, b) {
    a = m({children: void 0}, b);
    if (b = db(b.children))
      a.children = b;
    return a;
  }
  function fb(a, b, c, d) {
    a = a.options;
    if (b) {
      b = {};
      for (var e = 0; e < c.length; e++)
        b["$" + c[e]] = true;
      for (c = 0; c < a.length; c++)
        e = b.hasOwnProperty("$" + a[c].value), a[c].selected !== e && (a[c].selected = e), e && d && (a[c].defaultSelected = true);
    } else {
      c = "" + Sa(c);
      b = null;
      for (e = 0; e < a.length; e++) {
        if (a[e].value === c) {
          a[e].selected = true;
          d && (a[e].defaultSelected = true);
          return;
        }
        b !== null || a[e].disabled || (b = a[e]);
      }
      b !== null && (b.selected = true);
    }
  }
  function gb(a, b) {
    if (b.dangerouslySetInnerHTML != null)
      throw Error(y(91));
    return m({}, b, {value: void 0, defaultValue: void 0, children: "" + a._wrapperState.initialValue});
  }
  function hb(a, b) {
    var c = b.value;
    if (c == null) {
      c = b.children;
      b = b.defaultValue;
      if (c != null) {
        if (b != null)
          throw Error(y(92));
        if (Array.isArray(c)) {
          if (!(1 >= c.length))
            throw Error(y(93));
          c = c[0];
        }
        b = c;
      }
      b == null && (b = "");
      c = b;
    }
    a._wrapperState = {initialValue: Sa(c)};
  }
  function ib(a, b) {
    var c = Sa(b.value), d = Sa(b.defaultValue);
    c != null && (c = "" + c, c !== a.value && (a.value = c), b.defaultValue == null && a.defaultValue !== c && (a.defaultValue = c));
    d != null && (a.defaultValue = "" + d);
  }
  function jb(a) {
    var b = a.textContent;
    b === a._wrapperState.initialValue && b !== "" && b !== null && (a.value = b);
  }
  var kb = {html: "http://www.w3.org/1999/xhtml", mathml: "http://www.w3.org/1998/Math/MathML", svg: "http://www.w3.org/2000/svg"};
  function lb(a) {
    switch (a) {
      case "svg":
        return "http://www.w3.org/2000/svg";
      case "math":
        return "http://www.w3.org/1998/Math/MathML";
      default:
        return "http://www.w3.org/1999/xhtml";
    }
  }
  function mb(a, b) {
    return a == null || a === "http://www.w3.org/1999/xhtml" ? lb(b) : a === "http://www.w3.org/2000/svg" && b === "foreignObject" ? "http://www.w3.org/1999/xhtml" : a;
  }
  var nb;
  var ob = function(a) {
    return typeof MSApp !== "undefined" && MSApp.execUnsafeLocalFunction ? function(b, c, d, e) {
      MSApp.execUnsafeLocalFunction(function() {
        return a(b, c, d, e);
      });
    } : a;
  }(function(a, b) {
    if (a.namespaceURI !== kb.svg || "innerHTML" in a)
      a.innerHTML = b;
    else {
      nb = nb || document.createElement("div");
      nb.innerHTML = "<svg>" + b.valueOf().toString() + "</svg>";
      for (b = nb.firstChild; a.firstChild; )
        a.removeChild(a.firstChild);
      for (; b.firstChild; )
        a.appendChild(b.firstChild);
    }
  });
  function pb(a, b) {
    if (b) {
      var c = a.firstChild;
      if (c && c === a.lastChild && c.nodeType === 3) {
        c.nodeValue = b;
        return;
      }
    }
    a.textContent = b;
  }
  var qb = {
    animationIterationCount: true,
    borderImageOutset: true,
    borderImageSlice: true,
    borderImageWidth: true,
    boxFlex: true,
    boxFlexGroup: true,
    boxOrdinalGroup: true,
    columnCount: true,
    columns: true,
    flex: true,
    flexGrow: true,
    flexPositive: true,
    flexShrink: true,
    flexNegative: true,
    flexOrder: true,
    gridArea: true,
    gridRow: true,
    gridRowEnd: true,
    gridRowSpan: true,
    gridRowStart: true,
    gridColumn: true,
    gridColumnEnd: true,
    gridColumnSpan: true,
    gridColumnStart: true,
    fontWeight: true,
    lineClamp: true,
    lineHeight: true,
    opacity: true,
    order: true,
    orphans: true,
    tabSize: true,
    widows: true,
    zIndex: true,
    zoom: true,
    fillOpacity: true,
    floodOpacity: true,
    stopOpacity: true,
    strokeDasharray: true,
    strokeDashoffset: true,
    strokeMiterlimit: true,
    strokeOpacity: true,
    strokeWidth: true
  };
  var rb = ["Webkit", "ms", "Moz", "O"];
  Object.keys(qb).forEach(function(a) {
    rb.forEach(function(b) {
      b = b + a.charAt(0).toUpperCase() + a.substring(1);
      qb[b] = qb[a];
    });
  });
  function sb(a, b, c) {
    return b == null || typeof b === "boolean" || b === "" ? "" : c || typeof b !== "number" || b === 0 || qb.hasOwnProperty(a) && qb[a] ? ("" + b).trim() : b + "px";
  }
  function tb(a, b) {
    a = a.style;
    for (var c in b)
      if (b.hasOwnProperty(c)) {
        var d = c.indexOf("--") === 0, e = sb(c, b[c], d);
        c === "float" && (c = "cssFloat");
        d ? a.setProperty(c, e) : a[c] = e;
      }
  }
  var ub = m({menuitem: true}, {area: true, base: true, br: true, col: true, embed: true, hr: true, img: true, input: true, keygen: true, link: true, meta: true, param: true, source: true, track: true, wbr: true});
  function vb(a, b) {
    if (b) {
      if (ub[a] && (b.children != null || b.dangerouslySetInnerHTML != null))
        throw Error(y(137, a));
      if (b.dangerouslySetInnerHTML != null) {
        if (b.children != null)
          throw Error(y(60));
        if (!(typeof b.dangerouslySetInnerHTML === "object" && "__html" in b.dangerouslySetInnerHTML))
          throw Error(y(61));
      }
      if (b.style != null && typeof b.style !== "object")
        throw Error(y(62));
    }
  }
  function wb(a, b) {
    if (a.indexOf("-") === -1)
      return typeof b.is === "string";
    switch (a) {
      case "annotation-xml":
      case "color-profile":
      case "font-face":
      case "font-face-src":
      case "font-face-uri":
      case "font-face-format":
      case "font-face-name":
      case "missing-glyph":
        return false;
      default:
        return true;
    }
  }
  function xb(a) {
    a = a.target || a.srcElement || window;
    a.correspondingUseElement && (a = a.correspondingUseElement);
    return a.nodeType === 3 ? a.parentNode : a;
  }
  var yb = null;
  var zb = null;
  var Ab = null;
  function Bb(a) {
    if (a = Cb(a)) {
      if (typeof yb !== "function")
        throw Error(y(280));
      var b = a.stateNode;
      b && (b = Db(b), yb(a.stateNode, a.type, b));
    }
  }
  function Eb(a) {
    zb ? Ab ? Ab.push(a) : Ab = [a] : zb = a;
  }
  function Fb() {
    if (zb) {
      var a = zb, b = Ab;
      Ab = zb = null;
      Bb(a);
      if (b)
        for (a = 0; a < b.length; a++)
          Bb(b[a]);
    }
  }
  function Gb(a, b) {
    return a(b);
  }
  function Hb(a, b, c, d, e) {
    return a(b, c, d, e);
  }
  function Ib() {
  }
  var Jb = Gb;
  var Kb = false;
  var Lb = false;
  function Mb() {
    if (zb !== null || Ab !== null)
      Ib(), Fb();
  }
  function Nb(a, b, c) {
    if (Lb)
      return a(b, c);
    Lb = true;
    try {
      return Jb(a, b, c);
    } finally {
      Lb = false, Mb();
    }
  }
  function Ob(a, b) {
    var c = a.stateNode;
    if (c === null)
      return null;
    var d = Db(c);
    if (d === null)
      return null;
    c = d[b];
    a:
      switch (b) {
        case "onClick":
        case "onClickCapture":
        case "onDoubleClick":
        case "onDoubleClickCapture":
        case "onMouseDown":
        case "onMouseDownCapture":
        case "onMouseMove":
        case "onMouseMoveCapture":
        case "onMouseUp":
        case "onMouseUpCapture":
        case "onMouseEnter":
          (d = !d.disabled) || (a = a.type, d = !(a === "button" || a === "input" || a === "select" || a === "textarea"));
          a = !d;
          break a;
        default:
          a = false;
      }
    if (a)
      return null;
    if (c && typeof c !== "function")
      throw Error(y(231, b, typeof c));
    return c;
  }
  var Pb = false;
  if (fa)
    try {
      Qb = {};
      Object.defineProperty(Qb, "passive", {get: function() {
        Pb = true;
      }});
      window.addEventListener("test", Qb, Qb);
      window.removeEventListener("test", Qb, Qb);
    } catch (a) {
      Pb = false;
    }
  var Qb;
  function Rb(a, b, c, d, e, f, g, h, k) {
    var l = Array.prototype.slice.call(arguments, 3);
    try {
      b.apply(c, l);
    } catch (n) {
      this.onError(n);
    }
  }
  var Sb = false;
  var Tb = null;
  var Ub = false;
  var Vb = null;
  var Wb = {onError: function(a) {
    Sb = true;
    Tb = a;
  }};
  function Xb(a, b, c, d, e, f, g, h, k) {
    Sb = false;
    Tb = null;
    Rb.apply(Wb, arguments);
  }
  function Yb(a, b, c, d, e, f, g, h, k) {
    Xb.apply(this, arguments);
    if (Sb) {
      if (Sb) {
        var l = Tb;
        Sb = false;
        Tb = null;
      } else
        throw Error(y(198));
      Ub || (Ub = true, Vb = l);
    }
  }
  function Zb(a) {
    var b = a, c = a;
    if (a.alternate)
      for (; b.return; )
        b = b.return;
    else {
      a = b;
      do
        b = a, (b.flags & 1026) !== 0 && (c = b.return), a = b.return;
      while (a);
    }
    return b.tag === 3 ? c : null;
  }
  function $b(a) {
    if (a.tag === 13) {
      var b = a.memoizedState;
      b === null && (a = a.alternate, a !== null && (b = a.memoizedState));
      if (b !== null)
        return b.dehydrated;
    }
    return null;
  }
  function ac(a) {
    if (Zb(a) !== a)
      throw Error(y(188));
  }
  function bc(a) {
    var b = a.alternate;
    if (!b) {
      b = Zb(a);
      if (b === null)
        throw Error(y(188));
      return b !== a ? null : a;
    }
    for (var c = a, d = b; ; ) {
      var e = c.return;
      if (e === null)
        break;
      var f = e.alternate;
      if (f === null) {
        d = e.return;
        if (d !== null) {
          c = d;
          continue;
        }
        break;
      }
      if (e.child === f.child) {
        for (f = e.child; f; ) {
          if (f === c)
            return ac(e), a;
          if (f === d)
            return ac(e), b;
          f = f.sibling;
        }
        throw Error(y(188));
      }
      if (c.return !== d.return)
        c = e, d = f;
      else {
        for (var g = false, h = e.child; h; ) {
          if (h === c) {
            g = true;
            c = e;
            d = f;
            break;
          }
          if (h === d) {
            g = true;
            d = e;
            c = f;
            break;
          }
          h = h.sibling;
        }
        if (!g) {
          for (h = f.child; h; ) {
            if (h === c) {
              g = true;
              c = f;
              d = e;
              break;
            }
            if (h === d) {
              g = true;
              d = f;
              c = e;
              break;
            }
            h = h.sibling;
          }
          if (!g)
            throw Error(y(189));
        }
      }
      if (c.alternate !== d)
        throw Error(y(190));
    }
    if (c.tag !== 3)
      throw Error(y(188));
    return c.stateNode.current === c ? a : b;
  }
  function cc(a) {
    a = bc(a);
    if (!a)
      return null;
    for (var b = a; ; ) {
      if (b.tag === 5 || b.tag === 6)
        return b;
      if (b.child)
        b.child.return = b, b = b.child;
      else {
        if (b === a)
          break;
        for (; !b.sibling; ) {
          if (!b.return || b.return === a)
            return null;
          b = b.return;
        }
        b.sibling.return = b.return;
        b = b.sibling;
      }
    }
    return null;
  }
  function dc(a, b) {
    for (var c = a.alternate; b !== null; ) {
      if (b === a || b === c)
        return true;
      b = b.return;
    }
    return false;
  }
  var ec;
  var fc;
  var gc;
  var hc;
  var ic = false;
  var jc = [];
  var kc = null;
  var lc = null;
  var mc = null;
  var nc = new Map();
  var oc = new Map();
  var pc = [];
  var qc = "mousedown mouseup touchcancel touchend touchstart auxclick dblclick pointercancel pointerdown pointerup dragend dragstart drop compositionend compositionstart keydown keypress keyup input textInput copy cut paste click change contextmenu reset submit".split(" ");
  function rc(a, b, c, d, e) {
    return {blockedOn: a, domEventName: b, eventSystemFlags: c | 16, nativeEvent: e, targetContainers: [d]};
  }
  function sc(a, b) {
    switch (a) {
      case "focusin":
      case "focusout":
        kc = null;
        break;
      case "dragenter":
      case "dragleave":
        lc = null;
        break;
      case "mouseover":
      case "mouseout":
        mc = null;
        break;
      case "pointerover":
      case "pointerout":
        nc.delete(b.pointerId);
        break;
      case "gotpointercapture":
      case "lostpointercapture":
        oc.delete(b.pointerId);
    }
  }
  function tc(a, b, c, d, e, f) {
    if (a === null || a.nativeEvent !== f)
      return a = rc(b, c, d, e, f), b !== null && (b = Cb(b), b !== null && fc(b)), a;
    a.eventSystemFlags |= d;
    b = a.targetContainers;
    e !== null && b.indexOf(e) === -1 && b.push(e);
    return a;
  }
  function uc(a, b, c, d, e) {
    switch (b) {
      case "focusin":
        return kc = tc(kc, a, b, c, d, e), true;
      case "dragenter":
        return lc = tc(lc, a, b, c, d, e), true;
      case "mouseover":
        return mc = tc(mc, a, b, c, d, e), true;
      case "pointerover":
        var f = e.pointerId;
        nc.set(f, tc(nc.get(f) || null, a, b, c, d, e));
        return true;
      case "gotpointercapture":
        return f = e.pointerId, oc.set(f, tc(oc.get(f) || null, a, b, c, d, e)), true;
    }
    return false;
  }
  function vc(a) {
    var b = wc(a.target);
    if (b !== null) {
      var c = Zb(b);
      if (c !== null) {
        if (b = c.tag, b === 13) {
          if (b = $b(c), b !== null) {
            a.blockedOn = b;
            hc(a.lanePriority, function() {
              r.unstable_runWithPriority(a.priority, function() {
                gc(c);
              });
            });
            return;
          }
        } else if (b === 3 && c.stateNode.hydrate) {
          a.blockedOn = c.tag === 3 ? c.stateNode.containerInfo : null;
          return;
        }
      }
    }
    a.blockedOn = null;
  }
  function xc(a) {
    if (a.blockedOn !== null)
      return false;
    for (var b = a.targetContainers; 0 < b.length; ) {
      var c = yc(a.domEventName, a.eventSystemFlags, b[0], a.nativeEvent);
      if (c !== null)
        return b = Cb(c), b !== null && fc(b), a.blockedOn = c, false;
      b.shift();
    }
    return true;
  }
  function zc(a, b, c) {
    xc(a) && c.delete(b);
  }
  function Ac() {
    for (ic = false; 0 < jc.length; ) {
      var a = jc[0];
      if (a.blockedOn !== null) {
        a = Cb(a.blockedOn);
        a !== null && ec(a);
        break;
      }
      for (var b = a.targetContainers; 0 < b.length; ) {
        var c = yc(a.domEventName, a.eventSystemFlags, b[0], a.nativeEvent);
        if (c !== null) {
          a.blockedOn = c;
          break;
        }
        b.shift();
      }
      a.blockedOn === null && jc.shift();
    }
    kc !== null && xc(kc) && (kc = null);
    lc !== null && xc(lc) && (lc = null);
    mc !== null && xc(mc) && (mc = null);
    nc.forEach(zc);
    oc.forEach(zc);
  }
  function Bc(a, b) {
    a.blockedOn === b && (a.blockedOn = null, ic || (ic = true, r.unstable_scheduleCallback(r.unstable_NormalPriority, Ac)));
  }
  function Cc(a) {
    function b(b2) {
      return Bc(b2, a);
    }
    if (0 < jc.length) {
      Bc(jc[0], a);
      for (var c = 1; c < jc.length; c++) {
        var d = jc[c];
        d.blockedOn === a && (d.blockedOn = null);
      }
    }
    kc !== null && Bc(kc, a);
    lc !== null && Bc(lc, a);
    mc !== null && Bc(mc, a);
    nc.forEach(b);
    oc.forEach(b);
    for (c = 0; c < pc.length; c++)
      d = pc[c], d.blockedOn === a && (d.blockedOn = null);
    for (; 0 < pc.length && (c = pc[0], c.blockedOn === null); )
      vc(c), c.blockedOn === null && pc.shift();
  }
  function Dc(a, b) {
    var c = {};
    c[a.toLowerCase()] = b.toLowerCase();
    c["Webkit" + a] = "webkit" + b;
    c["Moz" + a] = "moz" + b;
    return c;
  }
  var Ec = {animationend: Dc("Animation", "AnimationEnd"), animationiteration: Dc("Animation", "AnimationIteration"), animationstart: Dc("Animation", "AnimationStart"), transitionend: Dc("Transition", "TransitionEnd")};
  var Fc = {};
  var Gc = {};
  fa && (Gc = document.createElement("div").style, "AnimationEvent" in window || (delete Ec.animationend.animation, delete Ec.animationiteration.animation, delete Ec.animationstart.animation), "TransitionEvent" in window || delete Ec.transitionend.transition);
  function Hc(a) {
    if (Fc[a])
      return Fc[a];
    if (!Ec[a])
      return a;
    var b = Ec[a], c;
    for (c in b)
      if (b.hasOwnProperty(c) && c in Gc)
        return Fc[a] = b[c];
    return a;
  }
  var Ic = Hc("animationend");
  var Jc = Hc("animationiteration");
  var Kc = Hc("animationstart");
  var Lc = Hc("transitionend");
  var Mc = new Map();
  var Nc = new Map();
  var Oc = [
    "abort",
    "abort",
    Ic,
    "animationEnd",
    Jc,
    "animationIteration",
    Kc,
    "animationStart",
    "canplay",
    "canPlay",
    "canplaythrough",
    "canPlayThrough",
    "durationchange",
    "durationChange",
    "emptied",
    "emptied",
    "encrypted",
    "encrypted",
    "ended",
    "ended",
    "error",
    "error",
    "gotpointercapture",
    "gotPointerCapture",
    "load",
    "load",
    "loadeddata",
    "loadedData",
    "loadedmetadata",
    "loadedMetadata",
    "loadstart",
    "loadStart",
    "lostpointercapture",
    "lostPointerCapture",
    "playing",
    "playing",
    "progress",
    "progress",
    "seeking",
    "seeking",
    "stalled",
    "stalled",
    "suspend",
    "suspend",
    "timeupdate",
    "timeUpdate",
    Lc,
    "transitionEnd",
    "waiting",
    "waiting"
  ];
  function Pc(a, b) {
    for (var c = 0; c < a.length; c += 2) {
      var d = a[c], e = a[c + 1];
      e = "on" + (e[0].toUpperCase() + e.slice(1));
      Nc.set(d, b);
      Mc.set(d, e);
      da(e, [d]);
    }
  }
  var Qc = r.unstable_now;
  Qc();
  var F = 8;
  function Rc(a) {
    if ((1 & a) !== 0)
      return F = 15, 1;
    if ((2 & a) !== 0)
      return F = 14, 2;
    if ((4 & a) !== 0)
      return F = 13, 4;
    var b = 24 & a;
    if (b !== 0)
      return F = 12, b;
    if ((a & 32) !== 0)
      return F = 11, 32;
    b = 192 & a;
    if (b !== 0)
      return F = 10, b;
    if ((a & 256) !== 0)
      return F = 9, 256;
    b = 3584 & a;
    if (b !== 0)
      return F = 8, b;
    if ((a & 4096) !== 0)
      return F = 7, 4096;
    b = 4186112 & a;
    if (b !== 0)
      return F = 6, b;
    b = 62914560 & a;
    if (b !== 0)
      return F = 5, b;
    if (a & 67108864)
      return F = 4, 67108864;
    if ((a & 134217728) !== 0)
      return F = 3, 134217728;
    b = 805306368 & a;
    if (b !== 0)
      return F = 2, b;
    if ((1073741824 & a) !== 0)
      return F = 1, 1073741824;
    F = 8;
    return a;
  }
  function Sc(a) {
    switch (a) {
      case 99:
        return 15;
      case 98:
        return 10;
      case 97:
      case 96:
        return 8;
      case 95:
        return 2;
      default:
        return 0;
    }
  }
  function Tc(a) {
    switch (a) {
      case 15:
      case 14:
        return 99;
      case 13:
      case 12:
      case 11:
      case 10:
        return 98;
      case 9:
      case 8:
      case 7:
      case 6:
      case 4:
      case 5:
        return 97;
      case 3:
      case 2:
      case 1:
        return 95;
      case 0:
        return 90;
      default:
        throw Error(y(358, a));
    }
  }
  function Uc(a, b) {
    var c = a.pendingLanes;
    if (c === 0)
      return F = 0;
    var d = 0, e = 0, f = a.expiredLanes, g = a.suspendedLanes, h = a.pingedLanes;
    if (f !== 0)
      d = f, e = F = 15;
    else if (f = c & 134217727, f !== 0) {
      var k = f & ~g;
      k !== 0 ? (d = Rc(k), e = F) : (h &= f, h !== 0 && (d = Rc(h), e = F));
    } else
      f = c & ~g, f !== 0 ? (d = Rc(f), e = F) : h !== 0 && (d = Rc(h), e = F);
    if (d === 0)
      return 0;
    d = 31 - Vc(d);
    d = c & ((0 > d ? 0 : 1 << d) << 1) - 1;
    if (b !== 0 && b !== d && (b & g) === 0) {
      Rc(b);
      if (e <= F)
        return b;
      F = e;
    }
    b = a.entangledLanes;
    if (b !== 0)
      for (a = a.entanglements, b &= d; 0 < b; )
        c = 31 - Vc(b), e = 1 << c, d |= a[c], b &= ~e;
    return d;
  }
  function Wc(a) {
    a = a.pendingLanes & -1073741825;
    return a !== 0 ? a : a & 1073741824 ? 1073741824 : 0;
  }
  function Xc(a, b) {
    switch (a) {
      case 15:
        return 1;
      case 14:
        return 2;
      case 12:
        return a = Yc(24 & ~b), a === 0 ? Xc(10, b) : a;
      case 10:
        return a = Yc(192 & ~b), a === 0 ? Xc(8, b) : a;
      case 8:
        return a = Yc(3584 & ~b), a === 0 && (a = Yc(4186112 & ~b), a === 0 && (a = 512)), a;
      case 2:
        return b = Yc(805306368 & ~b), b === 0 && (b = 268435456), b;
    }
    throw Error(y(358, a));
  }
  function Yc(a) {
    return a & -a;
  }
  function Zc(a) {
    for (var b = [], c = 0; 31 > c; c++)
      b.push(a);
    return b;
  }
  function $c(a, b, c) {
    a.pendingLanes |= b;
    var d = b - 1;
    a.suspendedLanes &= d;
    a.pingedLanes &= d;
    a = a.eventTimes;
    b = 31 - Vc(b);
    a[b] = c;
  }
  var Vc = Math.clz32 ? Math.clz32 : ad;
  var bd = Math.log;
  var cd = Math.LN2;
  function ad(a) {
    return a === 0 ? 32 : 31 - (bd(a) / cd | 0) | 0;
  }
  var dd = r.unstable_UserBlockingPriority;
  var ed = r.unstable_runWithPriority;
  var fd = true;
  function gd(a, b, c, d) {
    Kb || Ib();
    var e = hd, f = Kb;
    Kb = true;
    try {
      Hb(e, a, b, c, d);
    } finally {
      (Kb = f) || Mb();
    }
  }
  function id(a, b, c, d) {
    ed(dd, hd.bind(null, a, b, c, d));
  }
  function hd(a, b, c, d) {
    if (fd) {
      var e;
      if ((e = (b & 4) === 0) && 0 < jc.length && -1 < qc.indexOf(a))
        a = rc(null, a, b, c, d), jc.push(a);
      else {
        var f = yc(a, b, c, d);
        if (f === null)
          e && sc(a, d);
        else {
          if (e) {
            if (-1 < qc.indexOf(a)) {
              a = rc(f, a, b, c, d);
              jc.push(a);
              return;
            }
            if (uc(f, a, b, c, d))
              return;
            sc(a, d);
          }
          jd(a, b, d, null, c);
        }
      }
    }
  }
  function yc(a, b, c, d) {
    var e = xb(d);
    e = wc(e);
    if (e !== null) {
      var f = Zb(e);
      if (f === null)
        e = null;
      else {
        var g = f.tag;
        if (g === 13) {
          e = $b(f);
          if (e !== null)
            return e;
          e = null;
        } else if (g === 3) {
          if (f.stateNode.hydrate)
            return f.tag === 3 ? f.stateNode.containerInfo : null;
          e = null;
        } else
          f !== e && (e = null);
      }
    }
    jd(a, b, d, e, c);
    return null;
  }
  var kd = null;
  var ld = null;
  var md = null;
  function nd() {
    if (md)
      return md;
    var a, b = ld, c = b.length, d, e = "value" in kd ? kd.value : kd.textContent, f = e.length;
    for (a = 0; a < c && b[a] === e[a]; a++)
      ;
    var g = c - a;
    for (d = 1; d <= g && b[c - d] === e[f - d]; d++)
      ;
    return md = e.slice(a, 1 < d ? 1 - d : void 0);
  }
  function od(a) {
    var b = a.keyCode;
    "charCode" in a ? (a = a.charCode, a === 0 && b === 13 && (a = 13)) : a = b;
    a === 10 && (a = 13);
    return 32 <= a || a === 13 ? a : 0;
  }
  function pd() {
    return true;
  }
  function qd() {
    return false;
  }
  function rd(a) {
    function b(b2, d, e, f, g) {
      this._reactName = b2;
      this._targetInst = e;
      this.type = d;
      this.nativeEvent = f;
      this.target = g;
      this.currentTarget = null;
      for (var c in a)
        a.hasOwnProperty(c) && (b2 = a[c], this[c] = b2 ? b2(f) : f[c]);
      this.isDefaultPrevented = (f.defaultPrevented != null ? f.defaultPrevented : f.returnValue === false) ? pd : qd;
      this.isPropagationStopped = qd;
      return this;
    }
    m(b.prototype, {preventDefault: function() {
      this.defaultPrevented = true;
      var a2 = this.nativeEvent;
      a2 && (a2.preventDefault ? a2.preventDefault() : typeof a2.returnValue !== "unknown" && (a2.returnValue = false), this.isDefaultPrevented = pd);
    }, stopPropagation: function() {
      var a2 = this.nativeEvent;
      a2 && (a2.stopPropagation ? a2.stopPropagation() : typeof a2.cancelBubble !== "unknown" && (a2.cancelBubble = true), this.isPropagationStopped = pd);
    }, persist: function() {
    }, isPersistent: pd});
    return b;
  }
  var sd = {eventPhase: 0, bubbles: 0, cancelable: 0, timeStamp: function(a) {
    return a.timeStamp || Date.now();
  }, defaultPrevented: 0, isTrusted: 0};
  var td = rd(sd);
  var ud = m({}, sd, {view: 0, detail: 0});
  var vd = rd(ud);
  var wd;
  var xd;
  var yd;
  var Ad = m({}, ud, {screenX: 0, screenY: 0, clientX: 0, clientY: 0, pageX: 0, pageY: 0, ctrlKey: 0, shiftKey: 0, altKey: 0, metaKey: 0, getModifierState: zd, button: 0, buttons: 0, relatedTarget: function(a) {
    return a.relatedTarget === void 0 ? a.fromElement === a.srcElement ? a.toElement : a.fromElement : a.relatedTarget;
  }, movementX: function(a) {
    if ("movementX" in a)
      return a.movementX;
    a !== yd && (yd && a.type === "mousemove" ? (wd = a.screenX - yd.screenX, xd = a.screenY - yd.screenY) : xd = wd = 0, yd = a);
    return wd;
  }, movementY: function(a) {
    return "movementY" in a ? a.movementY : xd;
  }});
  var Bd = rd(Ad);
  var Cd = m({}, Ad, {dataTransfer: 0});
  var Dd = rd(Cd);
  var Ed = m({}, ud, {relatedTarget: 0});
  var Fd = rd(Ed);
  var Gd = m({}, sd, {animationName: 0, elapsedTime: 0, pseudoElement: 0});
  var Hd = rd(Gd);
  var Id = m({}, sd, {clipboardData: function(a) {
    return "clipboardData" in a ? a.clipboardData : window.clipboardData;
  }});
  var Jd = rd(Id);
  var Kd = m({}, sd, {data: 0});
  var Ld = rd(Kd);
  var Md = {
    Esc: "Escape",
    Spacebar: " ",
    Left: "ArrowLeft",
    Up: "ArrowUp",
    Right: "ArrowRight",
    Down: "ArrowDown",
    Del: "Delete",
    Win: "OS",
    Menu: "ContextMenu",
    Apps: "ContextMenu",
    Scroll: "ScrollLock",
    MozPrintableKey: "Unidentified"
  };
  var Nd = {
    8: "Backspace",
    9: "Tab",
    12: "Clear",
    13: "Enter",
    16: "Shift",
    17: "Control",
    18: "Alt",
    19: "Pause",
    20: "CapsLock",
    27: "Escape",
    32: " ",
    33: "PageUp",
    34: "PageDown",
    35: "End",
    36: "Home",
    37: "ArrowLeft",
    38: "ArrowUp",
    39: "ArrowRight",
    40: "ArrowDown",
    45: "Insert",
    46: "Delete",
    112: "F1",
    113: "F2",
    114: "F3",
    115: "F4",
    116: "F5",
    117: "F6",
    118: "F7",
    119: "F8",
    120: "F9",
    121: "F10",
    122: "F11",
    123: "F12",
    144: "NumLock",
    145: "ScrollLock",
    224: "Meta"
  };
  var Od = {Alt: "altKey", Control: "ctrlKey", Meta: "metaKey", Shift: "shiftKey"};
  function Pd(a) {
    var b = this.nativeEvent;
    return b.getModifierState ? b.getModifierState(a) : (a = Od[a]) ? !!b[a] : false;
  }
  function zd() {
    return Pd;
  }
  var Qd = m({}, ud, {key: function(a) {
    if (a.key) {
      var b = Md[a.key] || a.key;
      if (b !== "Unidentified")
        return b;
    }
    return a.type === "keypress" ? (a = od(a), a === 13 ? "Enter" : String.fromCharCode(a)) : a.type === "keydown" || a.type === "keyup" ? Nd[a.keyCode] || "Unidentified" : "";
  }, code: 0, location: 0, ctrlKey: 0, shiftKey: 0, altKey: 0, metaKey: 0, repeat: 0, locale: 0, getModifierState: zd, charCode: function(a) {
    return a.type === "keypress" ? od(a) : 0;
  }, keyCode: function(a) {
    return a.type === "keydown" || a.type === "keyup" ? a.keyCode : 0;
  }, which: function(a) {
    return a.type === "keypress" ? od(a) : a.type === "keydown" || a.type === "keyup" ? a.keyCode : 0;
  }});
  var Rd = rd(Qd);
  var Sd = m({}, Ad, {pointerId: 0, width: 0, height: 0, pressure: 0, tangentialPressure: 0, tiltX: 0, tiltY: 0, twist: 0, pointerType: 0, isPrimary: 0});
  var Td = rd(Sd);
  var Ud = m({}, ud, {touches: 0, targetTouches: 0, changedTouches: 0, altKey: 0, metaKey: 0, ctrlKey: 0, shiftKey: 0, getModifierState: zd});
  var Vd = rd(Ud);
  var Wd = m({}, sd, {propertyName: 0, elapsedTime: 0, pseudoElement: 0});
  var Xd = rd(Wd);
  var Yd = m({}, Ad, {
    deltaX: function(a) {
      return "deltaX" in a ? a.deltaX : "wheelDeltaX" in a ? -a.wheelDeltaX : 0;
    },
    deltaY: function(a) {
      return "deltaY" in a ? a.deltaY : "wheelDeltaY" in a ? -a.wheelDeltaY : "wheelDelta" in a ? -a.wheelDelta : 0;
    },
    deltaZ: 0,
    deltaMode: 0
  });
  var Zd = rd(Yd);
  var $d = [9, 13, 27, 32];
  var ae = fa && "CompositionEvent" in window;
  var be = null;
  fa && "documentMode" in document && (be = document.documentMode);
  var ce = fa && "TextEvent" in window && !be;
  var de = fa && (!ae || be && 8 < be && 11 >= be);
  var ee = String.fromCharCode(32);
  var fe = false;
  function ge(a, b) {
    switch (a) {
      case "keyup":
        return $d.indexOf(b.keyCode) !== -1;
      case "keydown":
        return b.keyCode !== 229;
      case "keypress":
      case "mousedown":
      case "focusout":
        return true;
      default:
        return false;
    }
  }
  function he(a) {
    a = a.detail;
    return typeof a === "object" && "data" in a ? a.data : null;
  }
  var ie = false;
  function je(a, b) {
    switch (a) {
      case "compositionend":
        return he(b);
      case "keypress":
        if (b.which !== 32)
          return null;
        fe = true;
        return ee;
      case "textInput":
        return a = b.data, a === ee && fe ? null : a;
      default:
        return null;
    }
  }
  function ke(a, b) {
    if (ie)
      return a === "compositionend" || !ae && ge(a, b) ? (a = nd(), md = ld = kd = null, ie = false, a) : null;
    switch (a) {
      case "paste":
        return null;
      case "keypress":
        if (!(b.ctrlKey || b.altKey || b.metaKey) || b.ctrlKey && b.altKey) {
          if (b.char && 1 < b.char.length)
            return b.char;
          if (b.which)
            return String.fromCharCode(b.which);
        }
        return null;
      case "compositionend":
        return de && b.locale !== "ko" ? null : b.data;
      default:
        return null;
    }
  }
  var le = {color: true, date: true, datetime: true, "datetime-local": true, email: true, month: true, number: true, password: true, range: true, search: true, tel: true, text: true, time: true, url: true, week: true};
  function me(a) {
    var b = a && a.nodeName && a.nodeName.toLowerCase();
    return b === "input" ? !!le[a.type] : b === "textarea" ? true : false;
  }
  function ne(a, b, c, d) {
    Eb(d);
    b = oe(b, "onChange");
    0 < b.length && (c = new td("onChange", "change", null, c, d), a.push({event: c, listeners: b}));
  }
  var pe = null;
  var qe = null;
  function re(a) {
    se(a, 0);
  }
  function te(a) {
    var b = ue(a);
    if (Wa(b))
      return a;
  }
  function ve(a, b) {
    if (a === "change")
      return b;
  }
  var we = false;
  if (fa) {
    if (fa) {
      ye = "oninput" in document;
      if (!ye) {
        ze = document.createElement("div");
        ze.setAttribute("oninput", "return;");
        ye = typeof ze.oninput === "function";
      }
      xe = ye;
    } else
      xe = false;
    we = xe && (!document.documentMode || 9 < document.documentMode);
  }
  var xe;
  var ye;
  var ze;
  function Ae() {
    pe && (pe.detachEvent("onpropertychange", Be), qe = pe = null);
  }
  function Be(a) {
    if (a.propertyName === "value" && te(qe)) {
      var b = [];
      ne(b, qe, a, xb(a));
      a = re;
      if (Kb)
        a(b);
      else {
        Kb = true;
        try {
          Gb(a, b);
        } finally {
          Kb = false, Mb();
        }
      }
    }
  }
  function Ce(a, b, c) {
    a === "focusin" ? (Ae(), pe = b, qe = c, pe.attachEvent("onpropertychange", Be)) : a === "focusout" && Ae();
  }
  function De(a) {
    if (a === "selectionchange" || a === "keyup" || a === "keydown")
      return te(qe);
  }
  function Ee(a, b) {
    if (a === "click")
      return te(b);
  }
  function Fe(a, b) {
    if (a === "input" || a === "change")
      return te(b);
  }
  function Ge(a, b) {
    return a === b && (a !== 0 || 1 / a === 1 / b) || a !== a && b !== b;
  }
  var He = typeof Object.is === "function" ? Object.is : Ge;
  var Ie = Object.prototype.hasOwnProperty;
  function Je(a, b) {
    if (He(a, b))
      return true;
    if (typeof a !== "object" || a === null || typeof b !== "object" || b === null)
      return false;
    var c = Object.keys(a), d = Object.keys(b);
    if (c.length !== d.length)
      return false;
    for (d = 0; d < c.length; d++)
      if (!Ie.call(b, c[d]) || !He(a[c[d]], b[c[d]]))
        return false;
    return true;
  }
  function Ke(a) {
    for (; a && a.firstChild; )
      a = a.firstChild;
    return a;
  }
  function Le(a, b) {
    var c = Ke(a);
    a = 0;
    for (var d; c; ) {
      if (c.nodeType === 3) {
        d = a + c.textContent.length;
        if (a <= b && d >= b)
          return {node: c, offset: b - a};
        a = d;
      }
      a: {
        for (; c; ) {
          if (c.nextSibling) {
            c = c.nextSibling;
            break a;
          }
          c = c.parentNode;
        }
        c = void 0;
      }
      c = Ke(c);
    }
  }
  function Me(a, b) {
    return a && b ? a === b ? true : a && a.nodeType === 3 ? false : b && b.nodeType === 3 ? Me(a, b.parentNode) : "contains" in a ? a.contains(b) : a.compareDocumentPosition ? !!(a.compareDocumentPosition(b) & 16) : false : false;
  }
  function Ne() {
    for (var a = window, b = Xa(); b instanceof a.HTMLIFrameElement; ) {
      try {
        var c = typeof b.contentWindow.location.href === "string";
      } catch (d) {
        c = false;
      }
      if (c)
        a = b.contentWindow;
      else
        break;
      b = Xa(a.document);
    }
    return b;
  }
  function Oe(a) {
    var b = a && a.nodeName && a.nodeName.toLowerCase();
    return b && (b === "input" && (a.type === "text" || a.type === "search" || a.type === "tel" || a.type === "url" || a.type === "password") || b === "textarea" || a.contentEditable === "true");
  }
  var Pe = fa && "documentMode" in document && 11 >= document.documentMode;
  var Qe = null;
  var Re = null;
  var Se = null;
  var Te = false;
  function Ue(a, b, c) {
    var d = c.window === c ? c.document : c.nodeType === 9 ? c : c.ownerDocument;
    Te || Qe == null || Qe !== Xa(d) || (d = Qe, "selectionStart" in d && Oe(d) ? d = {start: d.selectionStart, end: d.selectionEnd} : (d = (d.ownerDocument && d.ownerDocument.defaultView || window).getSelection(), d = {anchorNode: d.anchorNode, anchorOffset: d.anchorOffset, focusNode: d.focusNode, focusOffset: d.focusOffset}), Se && Je(Se, d) || (Se = d, d = oe(Re, "onSelect"), 0 < d.length && (b = new td("onSelect", "select", null, b, c), a.push({event: b, listeners: d}), b.target = Qe)));
  }
  Pc("cancel cancel click click close close contextmenu contextMenu copy copy cut cut auxclick auxClick dblclick doubleClick dragend dragEnd dragstart dragStart drop drop focusin focus focusout blur input input invalid invalid keydown keyDown keypress keyPress keyup keyUp mousedown mouseDown mouseup mouseUp paste paste pause pause play play pointercancel pointerCancel pointerdown pointerDown pointerup pointerUp ratechange rateChange reset reset seeked seeked submit submit touchcancel touchCancel touchend touchEnd touchstart touchStart volumechange volumeChange".split(" "), 0);
  Pc("drag drag dragenter dragEnter dragexit dragExit dragleave dragLeave dragover dragOver mousemove mouseMove mouseout mouseOut mouseover mouseOver pointermove pointerMove pointerout pointerOut pointerover pointerOver scroll scroll toggle toggle touchmove touchMove wheel wheel".split(" "), 1);
  Pc(Oc, 2);
  for (var Ve = "change selectionchange textInput compositionstart compositionend compositionupdate".split(" "), We = 0; We < Ve.length; We++)
    Nc.set(Ve[We], 0);
  ea("onMouseEnter", ["mouseout", "mouseover"]);
  ea("onMouseLeave", ["mouseout", "mouseover"]);
  ea("onPointerEnter", ["pointerout", "pointerover"]);
  ea("onPointerLeave", ["pointerout", "pointerover"]);
  da("onChange", "change click focusin focusout input keydown keyup selectionchange".split(" "));
  da("onSelect", "focusout contextmenu dragend focusin keydown keyup mousedown mouseup selectionchange".split(" "));
  da("onBeforeInput", ["compositionend", "keypress", "textInput", "paste"]);
  da("onCompositionEnd", "compositionend focusout keydown keypress keyup mousedown".split(" "));
  da("onCompositionStart", "compositionstart focusout keydown keypress keyup mousedown".split(" "));
  da("onCompositionUpdate", "compositionupdate focusout keydown keypress keyup mousedown".split(" "));
  var Xe = "abort canplay canplaythrough durationchange emptied encrypted ended error loadeddata loadedmetadata loadstart pause play playing progress ratechange seeked seeking stalled suspend timeupdate volumechange waiting".split(" ");
  var Ye = new Set("cancel close invalid load scroll toggle".split(" ").concat(Xe));
  function Ze(a, b, c) {
    var d = a.type || "unknown-event";
    a.currentTarget = c;
    Yb(d, b, void 0, a);
    a.currentTarget = null;
  }
  function se(a, b) {
    b = (b & 4) !== 0;
    for (var c = 0; c < a.length; c++) {
      var d = a[c], e = d.event;
      d = d.listeners;
      a: {
        var f = void 0;
        if (b)
          for (var g = d.length - 1; 0 <= g; g--) {
            var h = d[g], k = h.instance, l = h.currentTarget;
            h = h.listener;
            if (k !== f && e.isPropagationStopped())
              break a;
            Ze(e, h, l);
            f = k;
          }
        else
          for (g = 0; g < d.length; g++) {
            h = d[g];
            k = h.instance;
            l = h.currentTarget;
            h = h.listener;
            if (k !== f && e.isPropagationStopped())
              break a;
            Ze(e, h, l);
            f = k;
          }
      }
    }
    if (Ub)
      throw a = Vb, Ub = false, Vb = null, a;
  }
  function G(a, b) {
    var c = $e(b), d = a + "__bubble";
    c.has(d) || (af(b, a, 2, false), c.add(d));
  }
  var bf = "_reactListening" + Math.random().toString(36).slice(2);
  function cf(a) {
    a[bf] || (a[bf] = true, ba.forEach(function(b) {
      Ye.has(b) || df(b, false, a, null);
      df(b, true, a, null);
    }));
  }
  function df(a, b, c, d) {
    var e = 4 < arguments.length && arguments[4] !== void 0 ? arguments[4] : 0, f = c;
    a === "selectionchange" && c.nodeType !== 9 && (f = c.ownerDocument);
    if (d !== null && !b && Ye.has(a)) {
      if (a !== "scroll")
        return;
      e |= 2;
      f = d;
    }
    var g = $e(f), h = a + "__" + (b ? "capture" : "bubble");
    g.has(h) || (b && (e |= 4), af(f, a, e, b), g.add(h));
  }
  function af(a, b, c, d) {
    var e = Nc.get(b);
    switch (e === void 0 ? 2 : e) {
      case 0:
        e = gd;
        break;
      case 1:
        e = id;
        break;
      default:
        e = hd;
    }
    c = e.bind(null, b, c, a);
    e = void 0;
    !Pb || b !== "touchstart" && b !== "touchmove" && b !== "wheel" || (e = true);
    d ? e !== void 0 ? a.addEventListener(b, c, {capture: true, passive: e}) : a.addEventListener(b, c, true) : e !== void 0 ? a.addEventListener(b, c, {passive: e}) : a.addEventListener(b, c, false);
  }
  function jd(a, b, c, d, e) {
    var f = d;
    if ((b & 1) === 0 && (b & 2) === 0 && d !== null)
      a:
        for (; ; ) {
          if (d === null)
            return;
          var g = d.tag;
          if (g === 3 || g === 4) {
            var h = d.stateNode.containerInfo;
            if (h === e || h.nodeType === 8 && h.parentNode === e)
              break;
            if (g === 4)
              for (g = d.return; g !== null; ) {
                var k = g.tag;
                if (k === 3 || k === 4) {
                  if (k = g.stateNode.containerInfo, k === e || k.nodeType === 8 && k.parentNode === e)
                    return;
                }
                g = g.return;
              }
            for (; h !== null; ) {
              g = wc(h);
              if (g === null)
                return;
              k = g.tag;
              if (k === 5 || k === 6) {
                d = f = g;
                continue a;
              }
              h = h.parentNode;
            }
          }
          d = d.return;
        }
    Nb(function() {
      var d2 = f, e2 = xb(c), g2 = [];
      a: {
        var h2 = Mc.get(a);
        if (h2 !== void 0) {
          var k2 = td, x = a;
          switch (a) {
            case "keypress":
              if (od(c) === 0)
                break a;
            case "keydown":
            case "keyup":
              k2 = Rd;
              break;
            case "focusin":
              x = "focus";
              k2 = Fd;
              break;
            case "focusout":
              x = "blur";
              k2 = Fd;
              break;
            case "beforeblur":
            case "afterblur":
              k2 = Fd;
              break;
            case "click":
              if (c.button === 2)
                break a;
            case "auxclick":
            case "dblclick":
            case "mousedown":
            case "mousemove":
            case "mouseup":
            case "mouseout":
            case "mouseover":
            case "contextmenu":
              k2 = Bd;
              break;
            case "drag":
            case "dragend":
            case "dragenter":
            case "dragexit":
            case "dragleave":
            case "dragover":
            case "dragstart":
            case "drop":
              k2 = Dd;
              break;
            case "touchcancel":
            case "touchend":
            case "touchmove":
            case "touchstart":
              k2 = Vd;
              break;
            case Ic:
            case Jc:
            case Kc:
              k2 = Hd;
              break;
            case Lc:
              k2 = Xd;
              break;
            case "scroll":
              k2 = vd;
              break;
            case "wheel":
              k2 = Zd;
              break;
            case "copy":
            case "cut":
            case "paste":
              k2 = Jd;
              break;
            case "gotpointercapture":
            case "lostpointercapture":
            case "pointercancel":
            case "pointerdown":
            case "pointermove":
            case "pointerout":
            case "pointerover":
            case "pointerup":
              k2 = Td;
          }
          var w = (b & 4) !== 0, z = !w && a === "scroll", u = w ? h2 !== null ? h2 + "Capture" : null : h2;
          w = [];
          for (var t = d2, q; t !== null; ) {
            q = t;
            var v = q.stateNode;
            q.tag === 5 && v !== null && (q = v, u !== null && (v = Ob(t, u), v != null && w.push(ef(t, v, q))));
            if (z)
              break;
            t = t.return;
          }
          0 < w.length && (h2 = new k2(h2, x, null, c, e2), g2.push({event: h2, listeners: w}));
        }
      }
      if ((b & 7) === 0) {
        a: {
          h2 = a === "mouseover" || a === "pointerover";
          k2 = a === "mouseout" || a === "pointerout";
          if (h2 && (b & 16) === 0 && (x = c.relatedTarget || c.fromElement) && (wc(x) || x[ff]))
            break a;
          if (k2 || h2) {
            h2 = e2.window === e2 ? e2 : (h2 = e2.ownerDocument) ? h2.defaultView || h2.parentWindow : window;
            if (k2) {
              if (x = c.relatedTarget || c.toElement, k2 = d2, x = x ? wc(x) : null, x !== null && (z = Zb(x), x !== z || x.tag !== 5 && x.tag !== 6))
                x = null;
            } else
              k2 = null, x = d2;
            if (k2 !== x) {
              w = Bd;
              v = "onMouseLeave";
              u = "onMouseEnter";
              t = "mouse";
              if (a === "pointerout" || a === "pointerover")
                w = Td, v = "onPointerLeave", u = "onPointerEnter", t = "pointer";
              z = k2 == null ? h2 : ue(k2);
              q = x == null ? h2 : ue(x);
              h2 = new w(v, t + "leave", k2, c, e2);
              h2.target = z;
              h2.relatedTarget = q;
              v = null;
              wc(e2) === d2 && (w = new w(u, t + "enter", x, c, e2), w.target = q, w.relatedTarget = z, v = w);
              z = v;
              if (k2 && x)
                b: {
                  w = k2;
                  u = x;
                  t = 0;
                  for (q = w; q; q = gf(q))
                    t++;
                  q = 0;
                  for (v = u; v; v = gf(v))
                    q++;
                  for (; 0 < t - q; )
                    w = gf(w), t--;
                  for (; 0 < q - t; )
                    u = gf(u), q--;
                  for (; t--; ) {
                    if (w === u || u !== null && w === u.alternate)
                      break b;
                    w = gf(w);
                    u = gf(u);
                  }
                  w = null;
                }
              else
                w = null;
              k2 !== null && hf(g2, h2, k2, w, false);
              x !== null && z !== null && hf(g2, z, x, w, true);
            }
          }
        }
        a: {
          h2 = d2 ? ue(d2) : window;
          k2 = h2.nodeName && h2.nodeName.toLowerCase();
          if (k2 === "select" || k2 === "input" && h2.type === "file")
            var J = ve;
          else if (me(h2))
            if (we)
              J = Fe;
            else {
              J = De;
              var K = Ce;
            }
          else
            (k2 = h2.nodeName) && k2.toLowerCase() === "input" && (h2.type === "checkbox" || h2.type === "radio") && (J = Ee);
          if (J && (J = J(a, d2))) {
            ne(g2, J, c, e2);
            break a;
          }
          K && K(a, h2, d2);
          a === "focusout" && (K = h2._wrapperState) && K.controlled && h2.type === "number" && bb(h2, "number", h2.value);
        }
        K = d2 ? ue(d2) : window;
        switch (a) {
          case "focusin":
            if (me(K) || K.contentEditable === "true")
              Qe = K, Re = d2, Se = null;
            break;
          case "focusout":
            Se = Re = Qe = null;
            break;
          case "mousedown":
            Te = true;
            break;
          case "contextmenu":
          case "mouseup":
          case "dragend":
            Te = false;
            Ue(g2, c, e2);
            break;
          case "selectionchange":
            if (Pe)
              break;
          case "keydown":
          case "keyup":
            Ue(g2, c, e2);
        }
        var Q;
        if (ae)
          b: {
            switch (a) {
              case "compositionstart":
                var L = "onCompositionStart";
                break b;
              case "compositionend":
                L = "onCompositionEnd";
                break b;
              case "compositionupdate":
                L = "onCompositionUpdate";
                break b;
            }
            L = void 0;
          }
        else
          ie ? ge(a, c) && (L = "onCompositionEnd") : a === "keydown" && c.keyCode === 229 && (L = "onCompositionStart");
        L && (de && c.locale !== "ko" && (ie || L !== "onCompositionStart" ? L === "onCompositionEnd" && ie && (Q = nd()) : (kd = e2, ld = "value" in kd ? kd.value : kd.textContent, ie = true)), K = oe(d2, L), 0 < K.length && (L = new Ld(L, a, null, c, e2), g2.push({event: L, listeners: K}), Q ? L.data = Q : (Q = he(c), Q !== null && (L.data = Q))));
        if (Q = ce ? je(a, c) : ke(a, c))
          d2 = oe(d2, "onBeforeInput"), 0 < d2.length && (e2 = new Ld("onBeforeInput", "beforeinput", null, c, e2), g2.push({event: e2, listeners: d2}), e2.data = Q);
      }
      se(g2, b);
    });
  }
  function ef(a, b, c) {
    return {instance: a, listener: b, currentTarget: c};
  }
  function oe(a, b) {
    for (var c = b + "Capture", d = []; a !== null; ) {
      var e = a, f = e.stateNode;
      e.tag === 5 && f !== null && (e = f, f = Ob(a, c), f != null && d.unshift(ef(a, f, e)), f = Ob(a, b), f != null && d.push(ef(a, f, e)));
      a = a.return;
    }
    return d;
  }
  function gf(a) {
    if (a === null)
      return null;
    do
      a = a.return;
    while (a && a.tag !== 5);
    return a ? a : null;
  }
  function hf(a, b, c, d, e) {
    for (var f = b._reactName, g = []; c !== null && c !== d; ) {
      var h = c, k = h.alternate, l = h.stateNode;
      if (k !== null && k === d)
        break;
      h.tag === 5 && l !== null && (h = l, e ? (k = Ob(c, f), k != null && g.unshift(ef(c, k, h))) : e || (k = Ob(c, f), k != null && g.push(ef(c, k, h))));
      c = c.return;
    }
    g.length !== 0 && a.push({event: b, listeners: g});
  }
  function jf() {
  }
  var kf = null;
  var lf = null;
  function mf(a, b) {
    switch (a) {
      case "button":
      case "input":
      case "select":
      case "textarea":
        return !!b.autoFocus;
    }
    return false;
  }
  function nf(a, b) {
    return a === "textarea" || a === "option" || a === "noscript" || typeof b.children === "string" || typeof b.children === "number" || typeof b.dangerouslySetInnerHTML === "object" && b.dangerouslySetInnerHTML !== null && b.dangerouslySetInnerHTML.__html != null;
  }
  var of = typeof setTimeout === "function" ? setTimeout : void 0;
  var pf = typeof clearTimeout === "function" ? clearTimeout : void 0;
  function qf(a) {
    a.nodeType === 1 ? a.textContent = "" : a.nodeType === 9 && (a = a.body, a != null && (a.textContent = ""));
  }
  function rf(a) {
    for (; a != null; a = a.nextSibling) {
      var b = a.nodeType;
      if (b === 1 || b === 3)
        break;
    }
    return a;
  }
  function sf(a) {
    a = a.previousSibling;
    for (var b = 0; a; ) {
      if (a.nodeType === 8) {
        var c = a.data;
        if (c === "$" || c === "$!" || c === "$?") {
          if (b === 0)
            return a;
          b--;
        } else
          c === "/$" && b++;
      }
      a = a.previousSibling;
    }
    return null;
  }
  var tf = 0;
  function uf(a) {
    return {$$typeof: Ga, toString: a, valueOf: a};
  }
  var vf = Math.random().toString(36).slice(2);
  var wf = "__reactFiber$" + vf;
  var xf = "__reactProps$" + vf;
  var ff = "__reactContainer$" + vf;
  var yf = "__reactEvents$" + vf;
  function wc(a) {
    var b = a[wf];
    if (b)
      return b;
    for (var c = a.parentNode; c; ) {
      if (b = c[ff] || c[wf]) {
        c = b.alternate;
        if (b.child !== null || c !== null && c.child !== null)
          for (a = sf(a); a !== null; ) {
            if (c = a[wf])
              return c;
            a = sf(a);
          }
        return b;
      }
      a = c;
      c = a.parentNode;
    }
    return null;
  }
  function Cb(a) {
    a = a[wf] || a[ff];
    return !a || a.tag !== 5 && a.tag !== 6 && a.tag !== 13 && a.tag !== 3 ? null : a;
  }
  function ue(a) {
    if (a.tag === 5 || a.tag === 6)
      return a.stateNode;
    throw Error(y(33));
  }
  function Db(a) {
    return a[xf] || null;
  }
  function $e(a) {
    var b = a[yf];
    b === void 0 && (b = a[yf] = new Set());
    return b;
  }
  var zf = [];
  var Af = -1;
  function Bf(a) {
    return {current: a};
  }
  function H(a) {
    0 > Af || (a.current = zf[Af], zf[Af] = null, Af--);
  }
  function I(a, b) {
    Af++;
    zf[Af] = a.current;
    a.current = b;
  }
  var Cf = {};
  var M = Bf(Cf);
  var N = Bf(false);
  var Df = Cf;
  function Ef(a, b) {
    var c = a.type.contextTypes;
    if (!c)
      return Cf;
    var d = a.stateNode;
    if (d && d.__reactInternalMemoizedUnmaskedChildContext === b)
      return d.__reactInternalMemoizedMaskedChildContext;
    var e = {}, f;
    for (f in c)
      e[f] = b[f];
    d && (a = a.stateNode, a.__reactInternalMemoizedUnmaskedChildContext = b, a.__reactInternalMemoizedMaskedChildContext = e);
    return e;
  }
  function Ff(a) {
    a = a.childContextTypes;
    return a !== null && a !== void 0;
  }
  function Gf() {
    H(N);
    H(M);
  }
  function Hf(a, b, c) {
    if (M.current !== Cf)
      throw Error(y(168));
    I(M, b);
    I(N, c);
  }
  function If(a, b, c) {
    var d = a.stateNode;
    a = b.childContextTypes;
    if (typeof d.getChildContext !== "function")
      return c;
    d = d.getChildContext();
    for (var e in d)
      if (!(e in a))
        throw Error(y(108, Ra(b) || "Unknown", e));
    return m({}, c, d);
  }
  function Jf(a) {
    a = (a = a.stateNode) && a.__reactInternalMemoizedMergedChildContext || Cf;
    Df = M.current;
    I(M, a);
    I(N, N.current);
    return true;
  }
  function Kf(a, b, c) {
    var d = a.stateNode;
    if (!d)
      throw Error(y(169));
    c ? (a = If(a, b, Df), d.__reactInternalMemoizedMergedChildContext = a, H(N), H(M), I(M, a)) : H(N);
    I(N, c);
  }
  var Lf = null;
  var Mf = null;
  var Nf = r.unstable_runWithPriority;
  var Of = r.unstable_scheduleCallback;
  var Pf = r.unstable_cancelCallback;
  var Qf = r.unstable_shouldYield;
  var Rf = r.unstable_requestPaint;
  var Sf = r.unstable_now;
  var Tf = r.unstable_getCurrentPriorityLevel;
  var Uf = r.unstable_ImmediatePriority;
  var Vf = r.unstable_UserBlockingPriority;
  var Wf = r.unstable_NormalPriority;
  var Xf = r.unstable_LowPriority;
  var Yf = r.unstable_IdlePriority;
  var Zf = {};
  var $f = Rf !== void 0 ? Rf : function() {
  };
  var ag = null;
  var bg = null;
  var cg = false;
  var dg = Sf();
  var O = 1e4 > dg ? Sf : function() {
    return Sf() - dg;
  };
  function eg() {
    switch (Tf()) {
      case Uf:
        return 99;
      case Vf:
        return 98;
      case Wf:
        return 97;
      case Xf:
        return 96;
      case Yf:
        return 95;
      default:
        throw Error(y(332));
    }
  }
  function fg(a) {
    switch (a) {
      case 99:
        return Uf;
      case 98:
        return Vf;
      case 97:
        return Wf;
      case 96:
        return Xf;
      case 95:
        return Yf;
      default:
        throw Error(y(332));
    }
  }
  function gg(a, b) {
    a = fg(a);
    return Nf(a, b);
  }
  function hg(a, b, c) {
    a = fg(a);
    return Of(a, b, c);
  }
  function ig() {
    if (bg !== null) {
      var a = bg;
      bg = null;
      Pf(a);
    }
    jg();
  }
  function jg() {
    if (!cg && ag !== null) {
      cg = true;
      var a = 0;
      try {
        var b = ag;
        gg(99, function() {
          for (; a < b.length; a++) {
            var c = b[a];
            do
              c = c(true);
            while (c !== null);
          }
        });
        ag = null;
      } catch (c) {
        throw ag !== null && (ag = ag.slice(a + 1)), Of(Uf, ig), c;
      } finally {
        cg = false;
      }
    }
  }
  var kg = ra.ReactCurrentBatchConfig;
  function lg(a, b) {
    if (a && a.defaultProps) {
      b = m({}, b);
      a = a.defaultProps;
      for (var c in a)
        b[c] === void 0 && (b[c] = a[c]);
      return b;
    }
    return b;
  }
  var mg = Bf(null);
  var ng = null;
  var og = null;
  var pg = null;
  function qg() {
    pg = og = ng = null;
  }
  function rg(a) {
    var b = mg.current;
    H(mg);
    a.type._context._currentValue = b;
  }
  function sg(a, b) {
    for (; a !== null; ) {
      var c = a.alternate;
      if ((a.childLanes & b) === b)
        if (c === null || (c.childLanes & b) === b)
          break;
        else
          c.childLanes |= b;
      else
        a.childLanes |= b, c !== null && (c.childLanes |= b);
      a = a.return;
    }
  }
  function tg(a, b) {
    ng = a;
    pg = og = null;
    a = a.dependencies;
    a !== null && a.firstContext !== null && ((a.lanes & b) !== 0 && (ug = true), a.firstContext = null);
  }
  function vg(a, b) {
    if (pg !== a && b !== false && b !== 0) {
      if (typeof b !== "number" || b === 1073741823)
        pg = a, b = 1073741823;
      b = {context: a, observedBits: b, next: null};
      if (og === null) {
        if (ng === null)
          throw Error(y(308));
        og = b;
        ng.dependencies = {lanes: 0, firstContext: b, responders: null};
      } else
        og = og.next = b;
    }
    return a._currentValue;
  }
  var wg = false;
  function xg(a) {
    a.updateQueue = {baseState: a.memoizedState, firstBaseUpdate: null, lastBaseUpdate: null, shared: {pending: null}, effects: null};
  }
  function yg(a, b) {
    a = a.updateQueue;
    b.updateQueue === a && (b.updateQueue = {baseState: a.baseState, firstBaseUpdate: a.firstBaseUpdate, lastBaseUpdate: a.lastBaseUpdate, shared: a.shared, effects: a.effects});
  }
  function zg(a, b) {
    return {eventTime: a, lane: b, tag: 0, payload: null, callback: null, next: null};
  }
  function Ag(a, b) {
    a = a.updateQueue;
    if (a !== null) {
      a = a.shared;
      var c = a.pending;
      c === null ? b.next = b : (b.next = c.next, c.next = b);
      a.pending = b;
    }
  }
  function Bg(a, b) {
    var c = a.updateQueue, d = a.alternate;
    if (d !== null && (d = d.updateQueue, c === d)) {
      var e = null, f = null;
      c = c.firstBaseUpdate;
      if (c !== null) {
        do {
          var g = {eventTime: c.eventTime, lane: c.lane, tag: c.tag, payload: c.payload, callback: c.callback, next: null};
          f === null ? e = f = g : f = f.next = g;
          c = c.next;
        } while (c !== null);
        f === null ? e = f = b : f = f.next = b;
      } else
        e = f = b;
      c = {baseState: d.baseState, firstBaseUpdate: e, lastBaseUpdate: f, shared: d.shared, effects: d.effects};
      a.updateQueue = c;
      return;
    }
    a = c.lastBaseUpdate;
    a === null ? c.firstBaseUpdate = b : a.next = b;
    c.lastBaseUpdate = b;
  }
  function Cg(a, b, c, d) {
    var e = a.updateQueue;
    wg = false;
    var f = e.firstBaseUpdate, g = e.lastBaseUpdate, h = e.shared.pending;
    if (h !== null) {
      e.shared.pending = null;
      var k = h, l = k.next;
      k.next = null;
      g === null ? f = l : g.next = l;
      g = k;
      var n = a.alternate;
      if (n !== null) {
        n = n.updateQueue;
        var A = n.lastBaseUpdate;
        A !== g && (A === null ? n.firstBaseUpdate = l : A.next = l, n.lastBaseUpdate = k);
      }
    }
    if (f !== null) {
      A = e.baseState;
      g = 0;
      n = l = k = null;
      do {
        h = f.lane;
        var p = f.eventTime;
        if ((d & h) === h) {
          n !== null && (n = n.next = {
            eventTime: p,
            lane: 0,
            tag: f.tag,
            payload: f.payload,
            callback: f.callback,
            next: null
          });
          a: {
            var C = a, x = f;
            h = b;
            p = c;
            switch (x.tag) {
              case 1:
                C = x.payload;
                if (typeof C === "function") {
                  A = C.call(p, A, h);
                  break a;
                }
                A = C;
                break a;
              case 3:
                C.flags = C.flags & -4097 | 64;
              case 0:
                C = x.payload;
                h = typeof C === "function" ? C.call(p, A, h) : C;
                if (h === null || h === void 0)
                  break a;
                A = m({}, A, h);
                break a;
              case 2:
                wg = true;
            }
          }
          f.callback !== null && (a.flags |= 32, h = e.effects, h === null ? e.effects = [f] : h.push(f));
        } else
          p = {eventTime: p, lane: h, tag: f.tag, payload: f.payload, callback: f.callback, next: null}, n === null ? (l = n = p, k = A) : n = n.next = p, g |= h;
        f = f.next;
        if (f === null)
          if (h = e.shared.pending, h === null)
            break;
          else
            f = h.next, h.next = null, e.lastBaseUpdate = h, e.shared.pending = null;
      } while (1);
      n === null && (k = A);
      e.baseState = k;
      e.firstBaseUpdate = l;
      e.lastBaseUpdate = n;
      Dg |= g;
      a.lanes = g;
      a.memoizedState = A;
    }
  }
  function Eg(a, b, c) {
    a = b.effects;
    b.effects = null;
    if (a !== null)
      for (b = 0; b < a.length; b++) {
        var d = a[b], e = d.callback;
        if (e !== null) {
          d.callback = null;
          d = c;
          if (typeof e !== "function")
            throw Error(y(191, e));
          e.call(d);
        }
      }
  }
  var Fg = new aa.Component().refs;
  function Gg(a, b, c, d) {
    b = a.memoizedState;
    c = c(d, b);
    c = c === null || c === void 0 ? b : m({}, b, c);
    a.memoizedState = c;
    a.lanes === 0 && (a.updateQueue.baseState = c);
  }
  var Kg = {isMounted: function(a) {
    return (a = a._reactInternals) ? Zb(a) === a : false;
  }, enqueueSetState: function(a, b, c) {
    a = a._reactInternals;
    var d = Hg(), e = Ig(a), f = zg(d, e);
    f.payload = b;
    c !== void 0 && c !== null && (f.callback = c);
    Ag(a, f);
    Jg(a, e, d);
  }, enqueueReplaceState: function(a, b, c) {
    a = a._reactInternals;
    var d = Hg(), e = Ig(a), f = zg(d, e);
    f.tag = 1;
    f.payload = b;
    c !== void 0 && c !== null && (f.callback = c);
    Ag(a, f);
    Jg(a, e, d);
  }, enqueueForceUpdate: function(a, b) {
    a = a._reactInternals;
    var c = Hg(), d = Ig(a), e = zg(c, d);
    e.tag = 2;
    b !== void 0 && b !== null && (e.callback = b);
    Ag(a, e);
    Jg(a, d, c);
  }};
  function Lg(a, b, c, d, e, f, g) {
    a = a.stateNode;
    return typeof a.shouldComponentUpdate === "function" ? a.shouldComponentUpdate(d, f, g) : b.prototype && b.prototype.isPureReactComponent ? !Je(c, d) || !Je(e, f) : true;
  }
  function Mg(a, b, c) {
    var d = false, e = Cf;
    var f = b.contextType;
    typeof f === "object" && f !== null ? f = vg(f) : (e = Ff(b) ? Df : M.current, d = b.contextTypes, f = (d = d !== null && d !== void 0) ? Ef(a, e) : Cf);
    b = new b(c, f);
    a.memoizedState = b.state !== null && b.state !== void 0 ? b.state : null;
    b.updater = Kg;
    a.stateNode = b;
    b._reactInternals = a;
    d && (a = a.stateNode, a.__reactInternalMemoizedUnmaskedChildContext = e, a.__reactInternalMemoizedMaskedChildContext = f);
    return b;
  }
  function Ng(a, b, c, d) {
    a = b.state;
    typeof b.componentWillReceiveProps === "function" && b.componentWillReceiveProps(c, d);
    typeof b.UNSAFE_componentWillReceiveProps === "function" && b.UNSAFE_componentWillReceiveProps(c, d);
    b.state !== a && Kg.enqueueReplaceState(b, b.state, null);
  }
  function Og(a, b, c, d) {
    var e = a.stateNode;
    e.props = c;
    e.state = a.memoizedState;
    e.refs = Fg;
    xg(a);
    var f = b.contextType;
    typeof f === "object" && f !== null ? e.context = vg(f) : (f = Ff(b) ? Df : M.current, e.context = Ef(a, f));
    Cg(a, c, e, d);
    e.state = a.memoizedState;
    f = b.getDerivedStateFromProps;
    typeof f === "function" && (Gg(a, b, f, c), e.state = a.memoizedState);
    typeof b.getDerivedStateFromProps === "function" || typeof e.getSnapshotBeforeUpdate === "function" || typeof e.UNSAFE_componentWillMount !== "function" && typeof e.componentWillMount !== "function" || (b = e.state, typeof e.componentWillMount === "function" && e.componentWillMount(), typeof e.UNSAFE_componentWillMount === "function" && e.UNSAFE_componentWillMount(), b !== e.state && Kg.enqueueReplaceState(e, e.state, null), Cg(a, c, e, d), e.state = a.memoizedState);
    typeof e.componentDidMount === "function" && (a.flags |= 4);
  }
  var Pg = Array.isArray;
  function Qg(a, b, c) {
    a = c.ref;
    if (a !== null && typeof a !== "function" && typeof a !== "object") {
      if (c._owner) {
        c = c._owner;
        if (c) {
          if (c.tag !== 1)
            throw Error(y(309));
          var d = c.stateNode;
        }
        if (!d)
          throw Error(y(147, a));
        var e = "" + a;
        if (b !== null && b.ref !== null && typeof b.ref === "function" && b.ref._stringRef === e)
          return b.ref;
        b = function(a2) {
          var b2 = d.refs;
          b2 === Fg && (b2 = d.refs = {});
          a2 === null ? delete b2[e] : b2[e] = a2;
        };
        b._stringRef = e;
        return b;
      }
      if (typeof a !== "string")
        throw Error(y(284));
      if (!c._owner)
        throw Error(y(290, a));
    }
    return a;
  }
  function Rg(a, b) {
    if (a.type !== "textarea")
      throw Error(y(31, Object.prototype.toString.call(b) === "[object Object]" ? "object with keys {" + Object.keys(b).join(", ") + "}" : b));
  }
  function Sg(a) {
    function b(b2, c2) {
      if (a) {
        var d2 = b2.lastEffect;
        d2 !== null ? (d2.nextEffect = c2, b2.lastEffect = c2) : b2.firstEffect = b2.lastEffect = c2;
        c2.nextEffect = null;
        c2.flags = 8;
      }
    }
    function c(c2, d2) {
      if (!a)
        return null;
      for (; d2 !== null; )
        b(c2, d2), d2 = d2.sibling;
      return null;
    }
    function d(a2, b2) {
      for (a2 = new Map(); b2 !== null; )
        b2.key !== null ? a2.set(b2.key, b2) : a2.set(b2.index, b2), b2 = b2.sibling;
      return a2;
    }
    function e(a2, b2) {
      a2 = Tg(a2, b2);
      a2.index = 0;
      a2.sibling = null;
      return a2;
    }
    function f(b2, c2, d2) {
      b2.index = d2;
      if (!a)
        return c2;
      d2 = b2.alternate;
      if (d2 !== null)
        return d2 = d2.index, d2 < c2 ? (b2.flags = 2, c2) : d2;
      b2.flags = 2;
      return c2;
    }
    function g(b2) {
      a && b2.alternate === null && (b2.flags = 2);
      return b2;
    }
    function h(a2, b2, c2, d2) {
      if (b2 === null || b2.tag !== 6)
        return b2 = Ug(c2, a2.mode, d2), b2.return = a2, b2;
      b2 = e(b2, c2);
      b2.return = a2;
      return b2;
    }
    function k(a2, b2, c2, d2) {
      if (b2 !== null && b2.elementType === c2.type)
        return d2 = e(b2, c2.props), d2.ref = Qg(a2, b2, c2), d2.return = a2, d2;
      d2 = Vg(c2.type, c2.key, c2.props, null, a2.mode, d2);
      d2.ref = Qg(a2, b2, c2);
      d2.return = a2;
      return d2;
    }
    function l(a2, b2, c2, d2) {
      if (b2 === null || b2.tag !== 4 || b2.stateNode.containerInfo !== c2.containerInfo || b2.stateNode.implementation !== c2.implementation)
        return b2 = Wg(c2, a2.mode, d2), b2.return = a2, b2;
      b2 = e(b2, c2.children || []);
      b2.return = a2;
      return b2;
    }
    function n(a2, b2, c2, d2, f2) {
      if (b2 === null || b2.tag !== 7)
        return b2 = Xg(c2, a2.mode, d2, f2), b2.return = a2, b2;
      b2 = e(b2, c2);
      b2.return = a2;
      return b2;
    }
    function A(a2, b2, c2) {
      if (typeof b2 === "string" || typeof b2 === "number")
        return b2 = Ug("" + b2, a2.mode, c2), b2.return = a2, b2;
      if (typeof b2 === "object" && b2 !== null) {
        switch (b2.$$typeof) {
          case sa:
            return c2 = Vg(b2.type, b2.key, b2.props, null, a2.mode, c2), c2.ref = Qg(a2, null, b2), c2.return = a2, c2;
          case ta:
            return b2 = Wg(b2, a2.mode, c2), b2.return = a2, b2;
        }
        if (Pg(b2) || La(b2))
          return b2 = Xg(b2, a2.mode, c2, null), b2.return = a2, b2;
        Rg(a2, b2);
      }
      return null;
    }
    function p(a2, b2, c2, d2) {
      var e2 = b2 !== null ? b2.key : null;
      if (typeof c2 === "string" || typeof c2 === "number")
        return e2 !== null ? null : h(a2, b2, "" + c2, d2);
      if (typeof c2 === "object" && c2 !== null) {
        switch (c2.$$typeof) {
          case sa:
            return c2.key === e2 ? c2.type === ua ? n(a2, b2, c2.props.children, d2, e2) : k(a2, b2, c2, d2) : null;
          case ta:
            return c2.key === e2 ? l(a2, b2, c2, d2) : null;
        }
        if (Pg(c2) || La(c2))
          return e2 !== null ? null : n(a2, b2, c2, d2, null);
        Rg(a2, c2);
      }
      return null;
    }
    function C(a2, b2, c2, d2, e2) {
      if (typeof d2 === "string" || typeof d2 === "number")
        return a2 = a2.get(c2) || null, h(b2, a2, "" + d2, e2);
      if (typeof d2 === "object" && d2 !== null) {
        switch (d2.$$typeof) {
          case sa:
            return a2 = a2.get(d2.key === null ? c2 : d2.key) || null, d2.type === ua ? n(b2, a2, d2.props.children, e2, d2.key) : k(b2, a2, d2, e2);
          case ta:
            return a2 = a2.get(d2.key === null ? c2 : d2.key) || null, l(b2, a2, d2, e2);
        }
        if (Pg(d2) || La(d2))
          return a2 = a2.get(c2) || null, n(b2, a2, d2, e2, null);
        Rg(b2, d2);
      }
      return null;
    }
    function x(e2, g2, h2, k2) {
      for (var l2 = null, t = null, u = g2, z = g2 = 0, q = null; u !== null && z < h2.length; z++) {
        u.index > z ? (q = u, u = null) : q = u.sibling;
        var n2 = p(e2, u, h2[z], k2);
        if (n2 === null) {
          u === null && (u = q);
          break;
        }
        a && u && n2.alternate === null && b(e2, u);
        g2 = f(n2, g2, z);
        t === null ? l2 = n2 : t.sibling = n2;
        t = n2;
        u = q;
      }
      if (z === h2.length)
        return c(e2, u), l2;
      if (u === null) {
        for (; z < h2.length; z++)
          u = A(e2, h2[z], k2), u !== null && (g2 = f(u, g2, z), t === null ? l2 = u : t.sibling = u, t = u);
        return l2;
      }
      for (u = d(e2, u); z < h2.length; z++)
        q = C(u, e2, z, h2[z], k2), q !== null && (a && q.alternate !== null && u.delete(q.key === null ? z : q.key), g2 = f(q, g2, z), t === null ? l2 = q : t.sibling = q, t = q);
      a && u.forEach(function(a2) {
        return b(e2, a2);
      });
      return l2;
    }
    function w(e2, g2, h2, k2) {
      var l2 = La(h2);
      if (typeof l2 !== "function")
        throw Error(y(150));
      h2 = l2.call(h2);
      if (h2 == null)
        throw Error(y(151));
      for (var t = l2 = null, u = g2, z = g2 = 0, q = null, n2 = h2.next(); u !== null && !n2.done; z++, n2 = h2.next()) {
        u.index > z ? (q = u, u = null) : q = u.sibling;
        var w2 = p(e2, u, n2.value, k2);
        if (w2 === null) {
          u === null && (u = q);
          break;
        }
        a && u && w2.alternate === null && b(e2, u);
        g2 = f(w2, g2, z);
        t === null ? l2 = w2 : t.sibling = w2;
        t = w2;
        u = q;
      }
      if (n2.done)
        return c(e2, u), l2;
      if (u === null) {
        for (; !n2.done; z++, n2 = h2.next())
          n2 = A(e2, n2.value, k2), n2 !== null && (g2 = f(n2, g2, z), t === null ? l2 = n2 : t.sibling = n2, t = n2);
        return l2;
      }
      for (u = d(e2, u); !n2.done; z++, n2 = h2.next())
        n2 = C(u, e2, z, n2.value, k2), n2 !== null && (a && n2.alternate !== null && u.delete(n2.key === null ? z : n2.key), g2 = f(n2, g2, z), t === null ? l2 = n2 : t.sibling = n2, t = n2);
      a && u.forEach(function(a2) {
        return b(e2, a2);
      });
      return l2;
    }
    return function(a2, d2, f2, h2) {
      var k2 = typeof f2 === "object" && f2 !== null && f2.type === ua && f2.key === null;
      k2 && (f2 = f2.props.children);
      var l2 = typeof f2 === "object" && f2 !== null;
      if (l2)
        switch (f2.$$typeof) {
          case sa:
            a: {
              l2 = f2.key;
              for (k2 = d2; k2 !== null; ) {
                if (k2.key === l2) {
                  switch (k2.tag) {
                    case 7:
                      if (f2.type === ua) {
                        c(a2, k2.sibling);
                        d2 = e(k2, f2.props.children);
                        d2.return = a2;
                        a2 = d2;
                        break a;
                      }
                      break;
                    default:
                      if (k2.elementType === f2.type) {
                        c(a2, k2.sibling);
                        d2 = e(k2, f2.props);
                        d2.ref = Qg(a2, k2, f2);
                        d2.return = a2;
                        a2 = d2;
                        break a;
                      }
                  }
                  c(a2, k2);
                  break;
                } else
                  b(a2, k2);
                k2 = k2.sibling;
              }
              f2.type === ua ? (d2 = Xg(f2.props.children, a2.mode, h2, f2.key), d2.return = a2, a2 = d2) : (h2 = Vg(f2.type, f2.key, f2.props, null, a2.mode, h2), h2.ref = Qg(a2, d2, f2), h2.return = a2, a2 = h2);
            }
            return g(a2);
          case ta:
            a: {
              for (k2 = f2.key; d2 !== null; ) {
                if (d2.key === k2)
                  if (d2.tag === 4 && d2.stateNode.containerInfo === f2.containerInfo && d2.stateNode.implementation === f2.implementation) {
                    c(a2, d2.sibling);
                    d2 = e(d2, f2.children || []);
                    d2.return = a2;
                    a2 = d2;
                    break a;
                  } else {
                    c(a2, d2);
                    break;
                  }
                else
                  b(a2, d2);
                d2 = d2.sibling;
              }
              d2 = Wg(f2, a2.mode, h2);
              d2.return = a2;
              a2 = d2;
            }
            return g(a2);
        }
      if (typeof f2 === "string" || typeof f2 === "number")
        return f2 = "" + f2, d2 !== null && d2.tag === 6 ? (c(a2, d2.sibling), d2 = e(d2, f2), d2.return = a2, a2 = d2) : (c(a2, d2), d2 = Ug(f2, a2.mode, h2), d2.return = a2, a2 = d2), g(a2);
      if (Pg(f2))
        return x(a2, d2, f2, h2);
      if (La(f2))
        return w(a2, d2, f2, h2);
      l2 && Rg(a2, f2);
      if (typeof f2 === "undefined" && !k2)
        switch (a2.tag) {
          case 1:
          case 22:
          case 0:
          case 11:
          case 15:
            throw Error(y(152, Ra(a2.type) || "Component"));
        }
      return c(a2, d2);
    };
  }
  var Yg = Sg(true);
  var Zg = Sg(false);
  var $g = {};
  var ah = Bf($g);
  var bh = Bf($g);
  var ch = Bf($g);
  function dh(a) {
    if (a === $g)
      throw Error(y(174));
    return a;
  }
  function eh(a, b) {
    I(ch, b);
    I(bh, a);
    I(ah, $g);
    a = b.nodeType;
    switch (a) {
      case 9:
      case 11:
        b = (b = b.documentElement) ? b.namespaceURI : mb(null, "");
        break;
      default:
        a = a === 8 ? b.parentNode : b, b = a.namespaceURI || null, a = a.tagName, b = mb(b, a);
    }
    H(ah);
    I(ah, b);
  }
  function fh() {
    H(ah);
    H(bh);
    H(ch);
  }
  function gh(a) {
    dh(ch.current);
    var b = dh(ah.current);
    var c = mb(b, a.type);
    b !== c && (I(bh, a), I(ah, c));
  }
  function hh(a) {
    bh.current === a && (H(ah), H(bh));
  }
  var P = Bf(0);
  function ih(a) {
    for (var b = a; b !== null; ) {
      if (b.tag === 13) {
        var c = b.memoizedState;
        if (c !== null && (c = c.dehydrated, c === null || c.data === "$?" || c.data === "$!"))
          return b;
      } else if (b.tag === 19 && b.memoizedProps.revealOrder !== void 0) {
        if ((b.flags & 64) !== 0)
          return b;
      } else if (b.child !== null) {
        b.child.return = b;
        b = b.child;
        continue;
      }
      if (b === a)
        break;
      for (; b.sibling === null; ) {
        if (b.return === null || b.return === a)
          return null;
        b = b.return;
      }
      b.sibling.return = b.return;
      b = b.sibling;
    }
    return null;
  }
  var jh = null;
  var kh = null;
  var lh = false;
  function mh(a, b) {
    var c = nh(5, null, null, 0);
    c.elementType = "DELETED";
    c.type = "DELETED";
    c.stateNode = b;
    c.return = a;
    c.flags = 8;
    a.lastEffect !== null ? (a.lastEffect.nextEffect = c, a.lastEffect = c) : a.firstEffect = a.lastEffect = c;
  }
  function oh(a, b) {
    switch (a.tag) {
      case 5:
        var c = a.type;
        b = b.nodeType !== 1 || c.toLowerCase() !== b.nodeName.toLowerCase() ? null : b;
        return b !== null ? (a.stateNode = b, true) : false;
      case 6:
        return b = a.pendingProps === "" || b.nodeType !== 3 ? null : b, b !== null ? (a.stateNode = b, true) : false;
      case 13:
        return false;
      default:
        return false;
    }
  }
  function ph(a) {
    if (lh) {
      var b = kh;
      if (b) {
        var c = b;
        if (!oh(a, b)) {
          b = rf(c.nextSibling);
          if (!b || !oh(a, b)) {
            a.flags = a.flags & -1025 | 2;
            lh = false;
            jh = a;
            return;
          }
          mh(jh, c);
        }
        jh = a;
        kh = rf(b.firstChild);
      } else
        a.flags = a.flags & -1025 | 2, lh = false, jh = a;
    }
  }
  function qh(a) {
    for (a = a.return; a !== null && a.tag !== 5 && a.tag !== 3 && a.tag !== 13; )
      a = a.return;
    jh = a;
  }
  function rh(a) {
    if (a !== jh)
      return false;
    if (!lh)
      return qh(a), lh = true, false;
    var b = a.type;
    if (a.tag !== 5 || b !== "head" && b !== "body" && !nf(b, a.memoizedProps))
      for (b = kh; b; )
        mh(a, b), b = rf(b.nextSibling);
    qh(a);
    if (a.tag === 13) {
      a = a.memoizedState;
      a = a !== null ? a.dehydrated : null;
      if (!a)
        throw Error(y(317));
      a: {
        a = a.nextSibling;
        for (b = 0; a; ) {
          if (a.nodeType === 8) {
            var c = a.data;
            if (c === "/$") {
              if (b === 0) {
                kh = rf(a.nextSibling);
                break a;
              }
              b--;
            } else
              c !== "$" && c !== "$!" && c !== "$?" || b++;
          }
          a = a.nextSibling;
        }
        kh = null;
      }
    } else
      kh = jh ? rf(a.stateNode.nextSibling) : null;
    return true;
  }
  function sh() {
    kh = jh = null;
    lh = false;
  }
  var th = [];
  function uh() {
    for (var a = 0; a < th.length; a++)
      th[a]._workInProgressVersionPrimary = null;
    th.length = 0;
  }
  var vh = ra.ReactCurrentDispatcher;
  var wh = ra.ReactCurrentBatchConfig;
  var xh = 0;
  var R = null;
  var S = null;
  var T = null;
  var yh = false;
  var zh = false;
  function Ah() {
    throw Error(y(321));
  }
  function Bh(a, b) {
    if (b === null)
      return false;
    for (var c = 0; c < b.length && c < a.length; c++)
      if (!He(a[c], b[c]))
        return false;
    return true;
  }
  function Ch(a, b, c, d, e, f) {
    xh = f;
    R = b;
    b.memoizedState = null;
    b.updateQueue = null;
    b.lanes = 0;
    vh.current = a === null || a.memoizedState === null ? Dh : Eh;
    a = c(d, e);
    if (zh) {
      f = 0;
      do {
        zh = false;
        if (!(25 > f))
          throw Error(y(301));
        f += 1;
        T = S = null;
        b.updateQueue = null;
        vh.current = Fh;
        a = c(d, e);
      } while (zh);
    }
    vh.current = Gh;
    b = S !== null && S.next !== null;
    xh = 0;
    T = S = R = null;
    yh = false;
    if (b)
      throw Error(y(300));
    return a;
  }
  function Hh() {
    var a = {memoizedState: null, baseState: null, baseQueue: null, queue: null, next: null};
    T === null ? R.memoizedState = T = a : T = T.next = a;
    return T;
  }
  function Ih() {
    if (S === null) {
      var a = R.alternate;
      a = a !== null ? a.memoizedState : null;
    } else
      a = S.next;
    var b = T === null ? R.memoizedState : T.next;
    if (b !== null)
      T = b, S = a;
    else {
      if (a === null)
        throw Error(y(310));
      S = a;
      a = {memoizedState: S.memoizedState, baseState: S.baseState, baseQueue: S.baseQueue, queue: S.queue, next: null};
      T === null ? R.memoizedState = T = a : T = T.next = a;
    }
    return T;
  }
  function Jh(a, b) {
    return typeof b === "function" ? b(a) : b;
  }
  function Kh(a) {
    var b = Ih(), c = b.queue;
    if (c === null)
      throw Error(y(311));
    c.lastRenderedReducer = a;
    var d = S, e = d.baseQueue, f = c.pending;
    if (f !== null) {
      if (e !== null) {
        var g = e.next;
        e.next = f.next;
        f.next = g;
      }
      d.baseQueue = e = f;
      c.pending = null;
    }
    if (e !== null) {
      e = e.next;
      d = d.baseState;
      var h = g = f = null, k = e;
      do {
        var l = k.lane;
        if ((xh & l) === l)
          h !== null && (h = h.next = {lane: 0, action: k.action, eagerReducer: k.eagerReducer, eagerState: k.eagerState, next: null}), d = k.eagerReducer === a ? k.eagerState : a(d, k.action);
        else {
          var n = {
            lane: l,
            action: k.action,
            eagerReducer: k.eagerReducer,
            eagerState: k.eagerState,
            next: null
          };
          h === null ? (g = h = n, f = d) : h = h.next = n;
          R.lanes |= l;
          Dg |= l;
        }
        k = k.next;
      } while (k !== null && k !== e);
      h === null ? f = d : h.next = g;
      He(d, b.memoizedState) || (ug = true);
      b.memoizedState = d;
      b.baseState = f;
      b.baseQueue = h;
      c.lastRenderedState = d;
    }
    return [b.memoizedState, c.dispatch];
  }
  function Lh(a) {
    var b = Ih(), c = b.queue;
    if (c === null)
      throw Error(y(311));
    c.lastRenderedReducer = a;
    var d = c.dispatch, e = c.pending, f = b.memoizedState;
    if (e !== null) {
      c.pending = null;
      var g = e = e.next;
      do
        f = a(f, g.action), g = g.next;
      while (g !== e);
      He(f, b.memoizedState) || (ug = true);
      b.memoizedState = f;
      b.baseQueue === null && (b.baseState = f);
      c.lastRenderedState = f;
    }
    return [f, d];
  }
  function Mh(a, b, c) {
    var d = b._getVersion;
    d = d(b._source);
    var e = b._workInProgressVersionPrimary;
    if (e !== null)
      a = e === d;
    else if (a = a.mutableReadLanes, a = (xh & a) === a)
      b._workInProgressVersionPrimary = d, th.push(b);
    if (a)
      return c(b._source);
    th.push(b);
    throw Error(y(350));
  }
  function Nh(a, b, c, d) {
    var e = U;
    if (e === null)
      throw Error(y(349));
    var f = b._getVersion, g = f(b._source), h = vh.current, k = h.useState(function() {
      return Mh(e, b, c);
    }), l = k[1], n = k[0];
    k = T;
    var A = a.memoizedState, p = A.refs, C = p.getSnapshot, x = A.source;
    A = A.subscribe;
    var w = R;
    a.memoizedState = {refs: p, source: b, subscribe: d};
    h.useEffect(function() {
      p.getSnapshot = c;
      p.setSnapshot = l;
      var a2 = f(b._source);
      if (!He(g, a2)) {
        a2 = c(b._source);
        He(n, a2) || (l(a2), a2 = Ig(w), e.mutableReadLanes |= a2 & e.pendingLanes);
        a2 = e.mutableReadLanes;
        e.entangledLanes |= a2;
        for (var d2 = e.entanglements, h2 = a2; 0 < h2; ) {
          var k2 = 31 - Vc(h2), v = 1 << k2;
          d2[k2] |= a2;
          h2 &= ~v;
        }
      }
    }, [c, b, d]);
    h.useEffect(function() {
      return d(b._source, function() {
        var a2 = p.getSnapshot, c2 = p.setSnapshot;
        try {
          c2(a2(b._source));
          var d2 = Ig(w);
          e.mutableReadLanes |= d2 & e.pendingLanes;
        } catch (q) {
          c2(function() {
            throw q;
          });
        }
      });
    }, [b, d]);
    He(C, c) && He(x, b) && He(A, d) || (a = {pending: null, dispatch: null, lastRenderedReducer: Jh, lastRenderedState: n}, a.dispatch = l = Oh.bind(null, R, a), k.queue = a, k.baseQueue = null, n = Mh(e, b, c), k.memoizedState = k.baseState = n);
    return n;
  }
  function Ph(a, b, c) {
    var d = Ih();
    return Nh(d, a, b, c);
  }
  function Qh(a) {
    var b = Hh();
    typeof a === "function" && (a = a());
    b.memoizedState = b.baseState = a;
    a = b.queue = {pending: null, dispatch: null, lastRenderedReducer: Jh, lastRenderedState: a};
    a = a.dispatch = Oh.bind(null, R, a);
    return [b.memoizedState, a];
  }
  function Rh(a, b, c, d) {
    a = {tag: a, create: b, destroy: c, deps: d, next: null};
    b = R.updateQueue;
    b === null ? (b = {lastEffect: null}, R.updateQueue = b, b.lastEffect = a.next = a) : (c = b.lastEffect, c === null ? b.lastEffect = a.next = a : (d = c.next, c.next = a, a.next = d, b.lastEffect = a));
    return a;
  }
  function Sh(a) {
    var b = Hh();
    a = {current: a};
    return b.memoizedState = a;
  }
  function Th() {
    return Ih().memoizedState;
  }
  function Uh(a, b, c, d) {
    var e = Hh();
    R.flags |= a;
    e.memoizedState = Rh(1 | b, c, void 0, d === void 0 ? null : d);
  }
  function Vh(a, b, c, d) {
    var e = Ih();
    d = d === void 0 ? null : d;
    var f = void 0;
    if (S !== null) {
      var g = S.memoizedState;
      f = g.destroy;
      if (d !== null && Bh(d, g.deps)) {
        Rh(b, c, f, d);
        return;
      }
    }
    R.flags |= a;
    e.memoizedState = Rh(1 | b, c, f, d);
  }
  function Wh(a, b) {
    return Uh(516, 4, a, b);
  }
  function Xh(a, b) {
    return Vh(516, 4, a, b);
  }
  function Yh(a, b) {
    return Vh(4, 2, a, b);
  }
  function Zh(a, b) {
    if (typeof b === "function")
      return a = a(), b(a), function() {
        b(null);
      };
    if (b !== null && b !== void 0)
      return a = a(), b.current = a, function() {
        b.current = null;
      };
  }
  function $h(a, b, c) {
    c = c !== null && c !== void 0 ? c.concat([a]) : null;
    return Vh(4, 2, Zh.bind(null, b, a), c);
  }
  function ai() {
  }
  function bi(a, b) {
    var c = Ih();
    b = b === void 0 ? null : b;
    var d = c.memoizedState;
    if (d !== null && b !== null && Bh(b, d[1]))
      return d[0];
    c.memoizedState = [a, b];
    return a;
  }
  function ci(a, b) {
    var c = Ih();
    b = b === void 0 ? null : b;
    var d = c.memoizedState;
    if (d !== null && b !== null && Bh(b, d[1]))
      return d[0];
    a = a();
    c.memoizedState = [a, b];
    return a;
  }
  function di(a, b) {
    var c = eg();
    gg(98 > c ? 98 : c, function() {
      a(true);
    });
    gg(97 < c ? 97 : c, function() {
      var c2 = wh.transition;
      wh.transition = 1;
      try {
        a(false), b();
      } finally {
        wh.transition = c2;
      }
    });
  }
  function Oh(a, b, c) {
    var d = Hg(), e = Ig(a), f = {lane: e, action: c, eagerReducer: null, eagerState: null, next: null}, g = b.pending;
    g === null ? f.next = f : (f.next = g.next, g.next = f);
    b.pending = f;
    g = a.alternate;
    if (a === R || g !== null && g === R)
      zh = yh = true;
    else {
      if (a.lanes === 0 && (g === null || g.lanes === 0) && (g = b.lastRenderedReducer, g !== null))
        try {
          var h = b.lastRenderedState, k = g(h, c);
          f.eagerReducer = g;
          f.eagerState = k;
          if (He(k, h))
            return;
        } catch (l) {
        } finally {
        }
      Jg(a, e, d);
    }
  }
  var Gh = {readContext: vg, useCallback: Ah, useContext: Ah, useEffect: Ah, useImperativeHandle: Ah, useLayoutEffect: Ah, useMemo: Ah, useReducer: Ah, useRef: Ah, useState: Ah, useDebugValue: Ah, useDeferredValue: Ah, useTransition: Ah, useMutableSource: Ah, useOpaqueIdentifier: Ah, unstable_isNewReconciler: false};
  var Dh = {readContext: vg, useCallback: function(a, b) {
    Hh().memoizedState = [a, b === void 0 ? null : b];
    return a;
  }, useContext: vg, useEffect: Wh, useImperativeHandle: function(a, b, c) {
    c = c !== null && c !== void 0 ? c.concat([a]) : null;
    return Uh(4, 2, Zh.bind(null, b, a), c);
  }, useLayoutEffect: function(a, b) {
    return Uh(4, 2, a, b);
  }, useMemo: function(a, b) {
    var c = Hh();
    b = b === void 0 ? null : b;
    a = a();
    c.memoizedState = [a, b];
    return a;
  }, useReducer: function(a, b, c) {
    var d = Hh();
    b = c !== void 0 ? c(b) : b;
    d.memoizedState = d.baseState = b;
    a = d.queue = {pending: null, dispatch: null, lastRenderedReducer: a, lastRenderedState: b};
    a = a.dispatch = Oh.bind(null, R, a);
    return [d.memoizedState, a];
  }, useRef: Sh, useState: Qh, useDebugValue: ai, useDeferredValue: function(a) {
    var b = Qh(a), c = b[0], d = b[1];
    Wh(function() {
      var b2 = wh.transition;
      wh.transition = 1;
      try {
        d(a);
      } finally {
        wh.transition = b2;
      }
    }, [a]);
    return c;
  }, useTransition: function() {
    var a = Qh(false), b = a[0];
    a = di.bind(null, a[1]);
    Sh(a);
    return [a, b];
  }, useMutableSource: function(a, b, c) {
    var d = Hh();
    d.memoizedState = {refs: {getSnapshot: b, setSnapshot: null}, source: a, subscribe: c};
    return Nh(d, a, b, c);
  }, useOpaqueIdentifier: function() {
    if (lh) {
      var a = false, b = uf(function() {
        a || (a = true, c("r:" + (tf++).toString(36)));
        throw Error(y(355));
      }), c = Qh(b)[1];
      (R.mode & 2) === 0 && (R.flags |= 516, Rh(5, function() {
        c("r:" + (tf++).toString(36));
      }, void 0, null));
      return b;
    }
    b = "r:" + (tf++).toString(36);
    Qh(b);
    return b;
  }, unstable_isNewReconciler: false};
  var Eh = {readContext: vg, useCallback: bi, useContext: vg, useEffect: Xh, useImperativeHandle: $h, useLayoutEffect: Yh, useMemo: ci, useReducer: Kh, useRef: Th, useState: function() {
    return Kh(Jh);
  }, useDebugValue: ai, useDeferredValue: function(a) {
    var b = Kh(Jh), c = b[0], d = b[1];
    Xh(function() {
      var b2 = wh.transition;
      wh.transition = 1;
      try {
        d(a);
      } finally {
        wh.transition = b2;
      }
    }, [a]);
    return c;
  }, useTransition: function() {
    var a = Kh(Jh)[0];
    return [
      Th().current,
      a
    ];
  }, useMutableSource: Ph, useOpaqueIdentifier: function() {
    return Kh(Jh)[0];
  }, unstable_isNewReconciler: false};
  var Fh = {readContext: vg, useCallback: bi, useContext: vg, useEffect: Xh, useImperativeHandle: $h, useLayoutEffect: Yh, useMemo: ci, useReducer: Lh, useRef: Th, useState: function() {
    return Lh(Jh);
  }, useDebugValue: ai, useDeferredValue: function(a) {
    var b = Lh(Jh), c = b[0], d = b[1];
    Xh(function() {
      var b2 = wh.transition;
      wh.transition = 1;
      try {
        d(a);
      } finally {
        wh.transition = b2;
      }
    }, [a]);
    return c;
  }, useTransition: function() {
    var a = Lh(Jh)[0];
    return [
      Th().current,
      a
    ];
  }, useMutableSource: Ph, useOpaqueIdentifier: function() {
    return Lh(Jh)[0];
  }, unstable_isNewReconciler: false};
  var ei = ra.ReactCurrentOwner;
  var ug = false;
  function fi(a, b, c, d) {
    b.child = a === null ? Zg(b, null, c, d) : Yg(b, a.child, c, d);
  }
  function gi(a, b, c, d, e) {
    c = c.render;
    var f = b.ref;
    tg(b, e);
    d = Ch(a, b, c, d, f, e);
    if (a !== null && !ug)
      return b.updateQueue = a.updateQueue, b.flags &= -517, a.lanes &= ~e, hi(a, b, e);
    b.flags |= 1;
    fi(a, b, d, e);
    return b.child;
  }
  function ii(a, b, c, d, e, f) {
    if (a === null) {
      var g = c.type;
      if (typeof g === "function" && !ji(g) && g.defaultProps === void 0 && c.compare === null && c.defaultProps === void 0)
        return b.tag = 15, b.type = g, ki(a, b, g, d, e, f);
      a = Vg(c.type, null, d, b, b.mode, f);
      a.ref = b.ref;
      a.return = b;
      return b.child = a;
    }
    g = a.child;
    if ((e & f) === 0 && (e = g.memoizedProps, c = c.compare, c = c !== null ? c : Je, c(e, d) && a.ref === b.ref))
      return hi(a, b, f);
    b.flags |= 1;
    a = Tg(g, d);
    a.ref = b.ref;
    a.return = b;
    return b.child = a;
  }
  function ki(a, b, c, d, e, f) {
    if (a !== null && Je(a.memoizedProps, d) && a.ref === b.ref)
      if (ug = false, (f & e) !== 0)
        (a.flags & 16384) !== 0 && (ug = true);
      else
        return b.lanes = a.lanes, hi(a, b, f);
    return li(a, b, c, d, f);
  }
  function mi(a, b, c) {
    var d = b.pendingProps, e = d.children, f = a !== null ? a.memoizedState : null;
    if (d.mode === "hidden" || d.mode === "unstable-defer-without-hiding")
      if ((b.mode & 4) === 0)
        b.memoizedState = {baseLanes: 0}, ni(b, c);
      else if ((c & 1073741824) !== 0)
        b.memoizedState = {baseLanes: 0}, ni(b, f !== null ? f.baseLanes : c);
      else
        return a = f !== null ? f.baseLanes | c : c, b.lanes = b.childLanes = 1073741824, b.memoizedState = {baseLanes: a}, ni(b, a), null;
    else
      f !== null ? (d = f.baseLanes | c, b.memoizedState = null) : d = c, ni(b, d);
    fi(a, b, e, c);
    return b.child;
  }
  function oi(a, b) {
    var c = b.ref;
    if (a === null && c !== null || a !== null && a.ref !== c)
      b.flags |= 128;
  }
  function li(a, b, c, d, e) {
    var f = Ff(c) ? Df : M.current;
    f = Ef(b, f);
    tg(b, e);
    c = Ch(a, b, c, d, f, e);
    if (a !== null && !ug)
      return b.updateQueue = a.updateQueue, b.flags &= -517, a.lanes &= ~e, hi(a, b, e);
    b.flags |= 1;
    fi(a, b, c, e);
    return b.child;
  }
  function pi(a, b, c, d, e) {
    if (Ff(c)) {
      var f = true;
      Jf(b);
    } else
      f = false;
    tg(b, e);
    if (b.stateNode === null)
      a !== null && (a.alternate = null, b.alternate = null, b.flags |= 2), Mg(b, c, d), Og(b, c, d, e), d = true;
    else if (a === null) {
      var g = b.stateNode, h = b.memoizedProps;
      g.props = h;
      var k = g.context, l = c.contextType;
      typeof l === "object" && l !== null ? l = vg(l) : (l = Ff(c) ? Df : M.current, l = Ef(b, l));
      var n = c.getDerivedStateFromProps, A = typeof n === "function" || typeof g.getSnapshotBeforeUpdate === "function";
      A || typeof g.UNSAFE_componentWillReceiveProps !== "function" && typeof g.componentWillReceiveProps !== "function" || (h !== d || k !== l) && Ng(b, g, d, l);
      wg = false;
      var p = b.memoizedState;
      g.state = p;
      Cg(b, d, g, e);
      k = b.memoizedState;
      h !== d || p !== k || N.current || wg ? (typeof n === "function" && (Gg(b, c, n, d), k = b.memoizedState), (h = wg || Lg(b, c, h, d, p, k, l)) ? (A || typeof g.UNSAFE_componentWillMount !== "function" && typeof g.componentWillMount !== "function" || (typeof g.componentWillMount === "function" && g.componentWillMount(), typeof g.UNSAFE_componentWillMount === "function" && g.UNSAFE_componentWillMount()), typeof g.componentDidMount === "function" && (b.flags |= 4)) : (typeof g.componentDidMount === "function" && (b.flags |= 4), b.memoizedProps = d, b.memoizedState = k), g.props = d, g.state = k, g.context = l, d = h) : (typeof g.componentDidMount === "function" && (b.flags |= 4), d = false);
    } else {
      g = b.stateNode;
      yg(a, b);
      h = b.memoizedProps;
      l = b.type === b.elementType ? h : lg(b.type, h);
      g.props = l;
      A = b.pendingProps;
      p = g.context;
      k = c.contextType;
      typeof k === "object" && k !== null ? k = vg(k) : (k = Ff(c) ? Df : M.current, k = Ef(b, k));
      var C = c.getDerivedStateFromProps;
      (n = typeof C === "function" || typeof g.getSnapshotBeforeUpdate === "function") || typeof g.UNSAFE_componentWillReceiveProps !== "function" && typeof g.componentWillReceiveProps !== "function" || (h !== A || p !== k) && Ng(b, g, d, k);
      wg = false;
      p = b.memoizedState;
      g.state = p;
      Cg(b, d, g, e);
      var x = b.memoizedState;
      h !== A || p !== x || N.current || wg ? (typeof C === "function" && (Gg(b, c, C, d), x = b.memoizedState), (l = wg || Lg(b, c, l, d, p, x, k)) ? (n || typeof g.UNSAFE_componentWillUpdate !== "function" && typeof g.componentWillUpdate !== "function" || (typeof g.componentWillUpdate === "function" && g.componentWillUpdate(d, x, k), typeof g.UNSAFE_componentWillUpdate === "function" && g.UNSAFE_componentWillUpdate(d, x, k)), typeof g.componentDidUpdate === "function" && (b.flags |= 4), typeof g.getSnapshotBeforeUpdate === "function" && (b.flags |= 256)) : (typeof g.componentDidUpdate !== "function" || h === a.memoizedProps && p === a.memoizedState || (b.flags |= 4), typeof g.getSnapshotBeforeUpdate !== "function" || h === a.memoizedProps && p === a.memoizedState || (b.flags |= 256), b.memoizedProps = d, b.memoizedState = x), g.props = d, g.state = x, g.context = k, d = l) : (typeof g.componentDidUpdate !== "function" || h === a.memoizedProps && p === a.memoizedState || (b.flags |= 4), typeof g.getSnapshotBeforeUpdate !== "function" || h === a.memoizedProps && p === a.memoizedState || (b.flags |= 256), d = false);
    }
    return qi(a, b, c, d, f, e);
  }
  function qi(a, b, c, d, e, f) {
    oi(a, b);
    var g = (b.flags & 64) !== 0;
    if (!d && !g)
      return e && Kf(b, c, false), hi(a, b, f);
    d = b.stateNode;
    ei.current = b;
    var h = g && typeof c.getDerivedStateFromError !== "function" ? null : d.render();
    b.flags |= 1;
    a !== null && g ? (b.child = Yg(b, a.child, null, f), b.child = Yg(b, null, h, f)) : fi(a, b, h, f);
    b.memoizedState = d.state;
    e && Kf(b, c, true);
    return b.child;
  }
  function ri(a) {
    var b = a.stateNode;
    b.pendingContext ? Hf(a, b.pendingContext, b.pendingContext !== b.context) : b.context && Hf(a, b.context, false);
    eh(a, b.containerInfo);
  }
  var si = {dehydrated: null, retryLane: 0};
  function ti(a, b, c) {
    var d = b.pendingProps, e = P.current, f = false, g;
    (g = (b.flags & 64) !== 0) || (g = a !== null && a.memoizedState === null ? false : (e & 2) !== 0);
    g ? (f = true, b.flags &= -65) : a !== null && a.memoizedState === null || d.fallback === void 0 || d.unstable_avoidThisFallback === true || (e |= 1);
    I(P, e & 1);
    if (a === null) {
      d.fallback !== void 0 && ph(b);
      a = d.children;
      e = d.fallback;
      if (f)
        return a = ui(b, a, e, c), b.child.memoizedState = {baseLanes: c}, b.memoizedState = si, a;
      if (typeof d.unstable_expectedLoadTime === "number")
        return a = ui(b, a, e, c), b.child.memoizedState = {baseLanes: c}, b.memoizedState = si, b.lanes = 33554432, a;
      c = vi({mode: "visible", children: a}, b.mode, c, null);
      c.return = b;
      return b.child = c;
    }
    if (a.memoizedState !== null) {
      if (f)
        return d = wi(a, b, d.children, d.fallback, c), f = b.child, e = a.child.memoizedState, f.memoizedState = e === null ? {baseLanes: c} : {baseLanes: e.baseLanes | c}, f.childLanes = a.childLanes & ~c, b.memoizedState = si, d;
      c = xi(a, b, d.children, c);
      b.memoizedState = null;
      return c;
    }
    if (f)
      return d = wi(a, b, d.children, d.fallback, c), f = b.child, e = a.child.memoizedState, f.memoizedState = e === null ? {baseLanes: c} : {baseLanes: e.baseLanes | c}, f.childLanes = a.childLanes & ~c, b.memoizedState = si, d;
    c = xi(a, b, d.children, c);
    b.memoizedState = null;
    return c;
  }
  function ui(a, b, c, d) {
    var e = a.mode, f = a.child;
    b = {mode: "hidden", children: b};
    (e & 2) === 0 && f !== null ? (f.childLanes = 0, f.pendingProps = b) : f = vi(b, e, 0, null);
    c = Xg(c, e, d, null);
    f.return = a;
    c.return = a;
    f.sibling = c;
    a.child = f;
    return c;
  }
  function xi(a, b, c, d) {
    var e = a.child;
    a = e.sibling;
    c = Tg(e, {mode: "visible", children: c});
    (b.mode & 2) === 0 && (c.lanes = d);
    c.return = b;
    c.sibling = null;
    a !== null && (a.nextEffect = null, a.flags = 8, b.firstEffect = b.lastEffect = a);
    return b.child = c;
  }
  function wi(a, b, c, d, e) {
    var f = b.mode, g = a.child;
    a = g.sibling;
    var h = {mode: "hidden", children: c};
    (f & 2) === 0 && b.child !== g ? (c = b.child, c.childLanes = 0, c.pendingProps = h, g = c.lastEffect, g !== null ? (b.firstEffect = c.firstEffect, b.lastEffect = g, g.nextEffect = null) : b.firstEffect = b.lastEffect = null) : c = Tg(g, h);
    a !== null ? d = Tg(a, d) : (d = Xg(d, f, e, null), d.flags |= 2);
    d.return = b;
    c.return = b;
    c.sibling = d;
    b.child = c;
    return d;
  }
  function yi(a, b) {
    a.lanes |= b;
    var c = a.alternate;
    c !== null && (c.lanes |= b);
    sg(a.return, b);
  }
  function zi(a, b, c, d, e, f) {
    var g = a.memoizedState;
    g === null ? a.memoizedState = {isBackwards: b, rendering: null, renderingStartTime: 0, last: d, tail: c, tailMode: e, lastEffect: f} : (g.isBackwards = b, g.rendering = null, g.renderingStartTime = 0, g.last = d, g.tail = c, g.tailMode = e, g.lastEffect = f);
  }
  function Ai(a, b, c) {
    var d = b.pendingProps, e = d.revealOrder, f = d.tail;
    fi(a, b, d.children, c);
    d = P.current;
    if ((d & 2) !== 0)
      d = d & 1 | 2, b.flags |= 64;
    else {
      if (a !== null && (a.flags & 64) !== 0)
        a:
          for (a = b.child; a !== null; ) {
            if (a.tag === 13)
              a.memoizedState !== null && yi(a, c);
            else if (a.tag === 19)
              yi(a, c);
            else if (a.child !== null) {
              a.child.return = a;
              a = a.child;
              continue;
            }
            if (a === b)
              break a;
            for (; a.sibling === null; ) {
              if (a.return === null || a.return === b)
                break a;
              a = a.return;
            }
            a.sibling.return = a.return;
            a = a.sibling;
          }
      d &= 1;
    }
    I(P, d);
    if ((b.mode & 2) === 0)
      b.memoizedState = null;
    else
      switch (e) {
        case "forwards":
          c = b.child;
          for (e = null; c !== null; )
            a = c.alternate, a !== null && ih(a) === null && (e = c), c = c.sibling;
          c = e;
          c === null ? (e = b.child, b.child = null) : (e = c.sibling, c.sibling = null);
          zi(b, false, e, c, f, b.lastEffect);
          break;
        case "backwards":
          c = null;
          e = b.child;
          for (b.child = null; e !== null; ) {
            a = e.alternate;
            if (a !== null && ih(a) === null) {
              b.child = e;
              break;
            }
            a = e.sibling;
            e.sibling = c;
            c = e;
            e = a;
          }
          zi(b, true, c, null, f, b.lastEffect);
          break;
        case "together":
          zi(b, false, null, null, void 0, b.lastEffect);
          break;
        default:
          b.memoizedState = null;
      }
    return b.child;
  }
  function hi(a, b, c) {
    a !== null && (b.dependencies = a.dependencies);
    Dg |= b.lanes;
    if ((c & b.childLanes) !== 0) {
      if (a !== null && b.child !== a.child)
        throw Error(y(153));
      if (b.child !== null) {
        a = b.child;
        c = Tg(a, a.pendingProps);
        b.child = c;
        for (c.return = b; a.sibling !== null; )
          a = a.sibling, c = c.sibling = Tg(a, a.pendingProps), c.return = b;
        c.sibling = null;
      }
      return b.child;
    }
    return null;
  }
  var Bi;
  var Ci;
  var Di;
  var Ei;
  Bi = function(a, b) {
    for (var c = b.child; c !== null; ) {
      if (c.tag === 5 || c.tag === 6)
        a.appendChild(c.stateNode);
      else if (c.tag !== 4 && c.child !== null) {
        c.child.return = c;
        c = c.child;
        continue;
      }
      if (c === b)
        break;
      for (; c.sibling === null; ) {
        if (c.return === null || c.return === b)
          return;
        c = c.return;
      }
      c.sibling.return = c.return;
      c = c.sibling;
    }
  };
  Ci = function() {
  };
  Di = function(a, b, c, d) {
    var e = a.memoizedProps;
    if (e !== d) {
      a = b.stateNode;
      dh(ah.current);
      var f = null;
      switch (c) {
        case "input":
          e = Ya(a, e);
          d = Ya(a, d);
          f = [];
          break;
        case "option":
          e = eb(a, e);
          d = eb(a, d);
          f = [];
          break;
        case "select":
          e = m({}, e, {value: void 0});
          d = m({}, d, {value: void 0});
          f = [];
          break;
        case "textarea":
          e = gb(a, e);
          d = gb(a, d);
          f = [];
          break;
        default:
          typeof e.onClick !== "function" && typeof d.onClick === "function" && (a.onclick = jf);
      }
      vb(c, d);
      var g;
      c = null;
      for (l in e)
        if (!d.hasOwnProperty(l) && e.hasOwnProperty(l) && e[l] != null)
          if (l === "style") {
            var h = e[l];
            for (g in h)
              h.hasOwnProperty(g) && (c || (c = {}), c[g] = "");
          } else
            l !== "dangerouslySetInnerHTML" && l !== "children" && l !== "suppressContentEditableWarning" && l !== "suppressHydrationWarning" && l !== "autoFocus" && (ca.hasOwnProperty(l) ? f || (f = []) : (f = f || []).push(l, null));
      for (l in d) {
        var k = d[l];
        h = e != null ? e[l] : void 0;
        if (d.hasOwnProperty(l) && k !== h && (k != null || h != null))
          if (l === "style")
            if (h) {
              for (g in h)
                !h.hasOwnProperty(g) || k && k.hasOwnProperty(g) || (c || (c = {}), c[g] = "");
              for (g in k)
                k.hasOwnProperty(g) && h[g] !== k[g] && (c || (c = {}), c[g] = k[g]);
            } else
              c || (f || (f = []), f.push(l, c)), c = k;
          else
            l === "dangerouslySetInnerHTML" ? (k = k ? k.__html : void 0, h = h ? h.__html : void 0, k != null && h !== k && (f = f || []).push(l, k)) : l === "children" ? typeof k !== "string" && typeof k !== "number" || (f = f || []).push(l, "" + k) : l !== "suppressContentEditableWarning" && l !== "suppressHydrationWarning" && (ca.hasOwnProperty(l) ? (k != null && l === "onScroll" && G("scroll", a), f || h === k || (f = [])) : typeof k === "object" && k !== null && k.$$typeof === Ga ? k.toString() : (f = f || []).push(l, k));
      }
      c && (f = f || []).push("style", c);
      var l = f;
      if (b.updateQueue = l)
        b.flags |= 4;
    }
  };
  Ei = function(a, b, c, d) {
    c !== d && (b.flags |= 4);
  };
  function Fi(a, b) {
    if (!lh)
      switch (a.tailMode) {
        case "hidden":
          b = a.tail;
          for (var c = null; b !== null; )
            b.alternate !== null && (c = b), b = b.sibling;
          c === null ? a.tail = null : c.sibling = null;
          break;
        case "collapsed":
          c = a.tail;
          for (var d = null; c !== null; )
            c.alternate !== null && (d = c), c = c.sibling;
          d === null ? b || a.tail === null ? a.tail = null : a.tail.sibling = null : d.sibling = null;
      }
  }
  function Gi(a, b, c) {
    var d = b.pendingProps;
    switch (b.tag) {
      case 2:
      case 16:
      case 15:
      case 0:
      case 11:
      case 7:
      case 8:
      case 12:
      case 9:
      case 14:
        return null;
      case 1:
        return Ff(b.type) && Gf(), null;
      case 3:
        fh();
        H(N);
        H(M);
        uh();
        d = b.stateNode;
        d.pendingContext && (d.context = d.pendingContext, d.pendingContext = null);
        if (a === null || a.child === null)
          rh(b) ? b.flags |= 4 : d.hydrate || (b.flags |= 256);
        Ci(b);
        return null;
      case 5:
        hh(b);
        var e = dh(ch.current);
        c = b.type;
        if (a !== null && b.stateNode != null)
          Di(a, b, c, d, e), a.ref !== b.ref && (b.flags |= 128);
        else {
          if (!d) {
            if (b.stateNode === null)
              throw Error(y(166));
            return null;
          }
          a = dh(ah.current);
          if (rh(b)) {
            d = b.stateNode;
            c = b.type;
            var f = b.memoizedProps;
            d[wf] = b;
            d[xf] = f;
            switch (c) {
              case "dialog":
                G("cancel", d);
                G("close", d);
                break;
              case "iframe":
              case "object":
              case "embed":
                G("load", d);
                break;
              case "video":
              case "audio":
                for (a = 0; a < Xe.length; a++)
                  G(Xe[a], d);
                break;
              case "source":
                G("error", d);
                break;
              case "img":
              case "image":
              case "link":
                G("error", d);
                G("load", d);
                break;
              case "details":
                G("toggle", d);
                break;
              case "input":
                Za(d, f);
                G("invalid", d);
                break;
              case "select":
                d._wrapperState = {wasMultiple: !!f.multiple};
                G("invalid", d);
                break;
              case "textarea":
                hb(d, f), G("invalid", d);
            }
            vb(c, f);
            a = null;
            for (var g in f)
              f.hasOwnProperty(g) && (e = f[g], g === "children" ? typeof e === "string" ? d.textContent !== e && (a = ["children", e]) : typeof e === "number" && d.textContent !== "" + e && (a = ["children", "" + e]) : ca.hasOwnProperty(g) && e != null && g === "onScroll" && G("scroll", d));
            switch (c) {
              case "input":
                Va(d);
                cb(d, f, true);
                break;
              case "textarea":
                Va(d);
                jb(d);
                break;
              case "select":
              case "option":
                break;
              default:
                typeof f.onClick === "function" && (d.onclick = jf);
            }
            d = a;
            b.updateQueue = d;
            d !== null && (b.flags |= 4);
          } else {
            g = e.nodeType === 9 ? e : e.ownerDocument;
            a === kb.html && (a = lb(c));
            a === kb.html ? c === "script" ? (a = g.createElement("div"), a.innerHTML = "<script></script>", a = a.removeChild(a.firstChild)) : typeof d.is === "string" ? a = g.createElement(c, {is: d.is}) : (a = g.createElement(c), c === "select" && (g = a, d.multiple ? g.multiple = true : d.size && (g.size = d.size))) : a = g.createElementNS(a, c);
            a[wf] = b;
            a[xf] = d;
            Bi(a, b, false, false);
            b.stateNode = a;
            g = wb(c, d);
            switch (c) {
              case "dialog":
                G("cancel", a);
                G("close", a);
                e = d;
                break;
              case "iframe":
              case "object":
              case "embed":
                G("load", a);
                e = d;
                break;
              case "video":
              case "audio":
                for (e = 0; e < Xe.length; e++)
                  G(Xe[e], a);
                e = d;
                break;
              case "source":
                G("error", a);
                e = d;
                break;
              case "img":
              case "image":
              case "link":
                G("error", a);
                G("load", a);
                e = d;
                break;
              case "details":
                G("toggle", a);
                e = d;
                break;
              case "input":
                Za(a, d);
                e = Ya(a, d);
                G("invalid", a);
                break;
              case "option":
                e = eb(a, d);
                break;
              case "select":
                a._wrapperState = {wasMultiple: !!d.multiple};
                e = m({}, d, {value: void 0});
                G("invalid", a);
                break;
              case "textarea":
                hb(a, d);
                e = gb(a, d);
                G("invalid", a);
                break;
              default:
                e = d;
            }
            vb(c, e);
            var h = e;
            for (f in h)
              if (h.hasOwnProperty(f)) {
                var k = h[f];
                f === "style" ? tb(a, k) : f === "dangerouslySetInnerHTML" ? (k = k ? k.__html : void 0, k != null && ob(a, k)) : f === "children" ? typeof k === "string" ? (c !== "textarea" || k !== "") && pb(a, k) : typeof k === "number" && pb(a, "" + k) : f !== "suppressContentEditableWarning" && f !== "suppressHydrationWarning" && f !== "autoFocus" && (ca.hasOwnProperty(f) ? k != null && f === "onScroll" && G("scroll", a) : k != null && qa(a, f, k, g));
              }
            switch (c) {
              case "input":
                Va(a);
                cb(a, d, false);
                break;
              case "textarea":
                Va(a);
                jb(a);
                break;
              case "option":
                d.value != null && a.setAttribute("value", "" + Sa(d.value));
                break;
              case "select":
                a.multiple = !!d.multiple;
                f = d.value;
                f != null ? fb(a, !!d.multiple, f, false) : d.defaultValue != null && fb(a, !!d.multiple, d.defaultValue, true);
                break;
              default:
                typeof e.onClick === "function" && (a.onclick = jf);
            }
            mf(c, d) && (b.flags |= 4);
          }
          b.ref !== null && (b.flags |= 128);
        }
        return null;
      case 6:
        if (a && b.stateNode != null)
          Ei(a, b, a.memoizedProps, d);
        else {
          if (typeof d !== "string" && b.stateNode === null)
            throw Error(y(166));
          c = dh(ch.current);
          dh(ah.current);
          rh(b) ? (d = b.stateNode, c = b.memoizedProps, d[wf] = b, d.nodeValue !== c && (b.flags |= 4)) : (d = (c.nodeType === 9 ? c : c.ownerDocument).createTextNode(d), d[wf] = b, b.stateNode = d);
        }
        return null;
      case 13:
        H(P);
        d = b.memoizedState;
        if ((b.flags & 64) !== 0)
          return b.lanes = c, b;
        d = d !== null;
        c = false;
        a === null ? b.memoizedProps.fallback !== void 0 && rh(b) : c = a.memoizedState !== null;
        if (d && !c && (b.mode & 2) !== 0)
          if (a === null && b.memoizedProps.unstable_avoidThisFallback !== true || (P.current & 1) !== 0)
            V === 0 && (V = 3);
          else {
            if (V === 0 || V === 3)
              V = 4;
            U === null || (Dg & 134217727) === 0 && (Hi & 134217727) === 0 || Ii(U, W);
          }
        if (d || c)
          b.flags |= 4;
        return null;
      case 4:
        return fh(), Ci(b), a === null && cf(b.stateNode.containerInfo), null;
      case 10:
        return rg(b), null;
      case 17:
        return Ff(b.type) && Gf(), null;
      case 19:
        H(P);
        d = b.memoizedState;
        if (d === null)
          return null;
        f = (b.flags & 64) !== 0;
        g = d.rendering;
        if (g === null)
          if (f)
            Fi(d, false);
          else {
            if (V !== 0 || a !== null && (a.flags & 64) !== 0)
              for (a = b.child; a !== null; ) {
                g = ih(a);
                if (g !== null) {
                  b.flags |= 64;
                  Fi(d, false);
                  f = g.updateQueue;
                  f !== null && (b.updateQueue = f, b.flags |= 4);
                  d.lastEffect === null && (b.firstEffect = null);
                  b.lastEffect = d.lastEffect;
                  d = c;
                  for (c = b.child; c !== null; )
                    f = c, a = d, f.flags &= 2, f.nextEffect = null, f.firstEffect = null, f.lastEffect = null, g = f.alternate, g === null ? (f.childLanes = 0, f.lanes = a, f.child = null, f.memoizedProps = null, f.memoizedState = null, f.updateQueue = null, f.dependencies = null, f.stateNode = null) : (f.childLanes = g.childLanes, f.lanes = g.lanes, f.child = g.child, f.memoizedProps = g.memoizedProps, f.memoizedState = g.memoizedState, f.updateQueue = g.updateQueue, f.type = g.type, a = g.dependencies, f.dependencies = a === null ? null : {lanes: a.lanes, firstContext: a.firstContext}), c = c.sibling;
                  I(P, P.current & 1 | 2);
                  return b.child;
                }
                a = a.sibling;
              }
            d.tail !== null && O() > Ji && (b.flags |= 64, f = true, Fi(d, false), b.lanes = 33554432);
          }
        else {
          if (!f)
            if (a = ih(g), a !== null) {
              if (b.flags |= 64, f = true, c = a.updateQueue, c !== null && (b.updateQueue = c, b.flags |= 4), Fi(d, true), d.tail === null && d.tailMode === "hidden" && !g.alternate && !lh)
                return b = b.lastEffect = d.lastEffect, b !== null && (b.nextEffect = null), null;
            } else
              2 * O() - d.renderingStartTime > Ji && c !== 1073741824 && (b.flags |= 64, f = true, Fi(d, false), b.lanes = 33554432);
          d.isBackwards ? (g.sibling = b.child, b.child = g) : (c = d.last, c !== null ? c.sibling = g : b.child = g, d.last = g);
        }
        return d.tail !== null ? (c = d.tail, d.rendering = c, d.tail = c.sibling, d.lastEffect = b.lastEffect, d.renderingStartTime = O(), c.sibling = null, b = P.current, I(P, f ? b & 1 | 2 : b & 1), c) : null;
      case 23:
      case 24:
        return Ki(), a !== null && a.memoizedState !== null !== (b.memoizedState !== null) && d.mode !== "unstable-defer-without-hiding" && (b.flags |= 4), null;
    }
    throw Error(y(156, b.tag));
  }
  function Li(a) {
    switch (a.tag) {
      case 1:
        Ff(a.type) && Gf();
        var b = a.flags;
        return b & 4096 ? (a.flags = b & -4097 | 64, a) : null;
      case 3:
        fh();
        H(N);
        H(M);
        uh();
        b = a.flags;
        if ((b & 64) !== 0)
          throw Error(y(285));
        a.flags = b & -4097 | 64;
        return a;
      case 5:
        return hh(a), null;
      case 13:
        return H(P), b = a.flags, b & 4096 ? (a.flags = b & -4097 | 64, a) : null;
      case 19:
        return H(P), null;
      case 4:
        return fh(), null;
      case 10:
        return rg(a), null;
      case 23:
      case 24:
        return Ki(), null;
      default:
        return null;
    }
  }
  function Mi(a, b) {
    try {
      var c = "", d = b;
      do
        c += Qa(d), d = d.return;
      while (d);
      var e = c;
    } catch (f) {
      e = "\nError generating stack: " + f.message + "\n" + f.stack;
    }
    return {value: a, source: b, stack: e};
  }
  function Ni(a, b) {
    try {
      console.error(b.value);
    } catch (c) {
      setTimeout(function() {
        throw c;
      });
    }
  }
  var Oi = typeof WeakMap === "function" ? WeakMap : Map;
  function Pi(a, b, c) {
    c = zg(-1, c);
    c.tag = 3;
    c.payload = {element: null};
    var d = b.value;
    c.callback = function() {
      Qi || (Qi = true, Ri = d);
      Ni(a, b);
    };
    return c;
  }
  function Si(a, b, c) {
    c = zg(-1, c);
    c.tag = 3;
    var d = a.type.getDerivedStateFromError;
    if (typeof d === "function") {
      var e = b.value;
      c.payload = function() {
        Ni(a, b);
        return d(e);
      };
    }
    var f = a.stateNode;
    f !== null && typeof f.componentDidCatch === "function" && (c.callback = function() {
      typeof d !== "function" && (Ti === null ? Ti = new Set([this]) : Ti.add(this), Ni(a, b));
      var c2 = b.stack;
      this.componentDidCatch(b.value, {componentStack: c2 !== null ? c2 : ""});
    });
    return c;
  }
  var Ui = typeof WeakSet === "function" ? WeakSet : Set;
  function Vi(a) {
    var b = a.ref;
    if (b !== null)
      if (typeof b === "function")
        try {
          b(null);
        } catch (c) {
          Wi(a, c);
        }
      else
        b.current = null;
  }
  function Xi(a, b) {
    switch (b.tag) {
      case 0:
      case 11:
      case 15:
      case 22:
        return;
      case 1:
        if (b.flags & 256 && a !== null) {
          var c = a.memoizedProps, d = a.memoizedState;
          a = b.stateNode;
          b = a.getSnapshotBeforeUpdate(b.elementType === b.type ? c : lg(b.type, c), d);
          a.__reactInternalSnapshotBeforeUpdate = b;
        }
        return;
      case 3:
        b.flags & 256 && qf(b.stateNode.containerInfo);
        return;
      case 5:
      case 6:
      case 4:
      case 17:
        return;
    }
    throw Error(y(163));
  }
  function Yi(a, b, c) {
    switch (c.tag) {
      case 0:
      case 11:
      case 15:
      case 22:
        b = c.updateQueue;
        b = b !== null ? b.lastEffect : null;
        if (b !== null) {
          a = b = b.next;
          do {
            if ((a.tag & 3) === 3) {
              var d = a.create;
              a.destroy = d();
            }
            a = a.next;
          } while (a !== b);
        }
        b = c.updateQueue;
        b = b !== null ? b.lastEffect : null;
        if (b !== null) {
          a = b = b.next;
          do {
            var e = a;
            d = e.next;
            e = e.tag;
            (e & 4) !== 0 && (e & 1) !== 0 && (Zi(c, a), $i(c, a));
            a = d;
          } while (a !== b);
        }
        return;
      case 1:
        a = c.stateNode;
        c.flags & 4 && (b === null ? a.componentDidMount() : (d = c.elementType === c.type ? b.memoizedProps : lg(c.type, b.memoizedProps), a.componentDidUpdate(d, b.memoizedState, a.__reactInternalSnapshotBeforeUpdate)));
        b = c.updateQueue;
        b !== null && Eg(c, b, a);
        return;
      case 3:
        b = c.updateQueue;
        if (b !== null) {
          a = null;
          if (c.child !== null)
            switch (c.child.tag) {
              case 5:
                a = c.child.stateNode;
                break;
              case 1:
                a = c.child.stateNode;
            }
          Eg(c, b, a);
        }
        return;
      case 5:
        a = c.stateNode;
        b === null && c.flags & 4 && mf(c.type, c.memoizedProps) && a.focus();
        return;
      case 6:
        return;
      case 4:
        return;
      case 12:
        return;
      case 13:
        c.memoizedState === null && (c = c.alternate, c !== null && (c = c.memoizedState, c !== null && (c = c.dehydrated, c !== null && Cc(c))));
        return;
      case 19:
      case 17:
      case 20:
      case 21:
      case 23:
      case 24:
        return;
    }
    throw Error(y(163));
  }
  function aj(a, b) {
    for (var c = a; ; ) {
      if (c.tag === 5) {
        var d = c.stateNode;
        if (b)
          d = d.style, typeof d.setProperty === "function" ? d.setProperty("display", "none", "important") : d.display = "none";
        else {
          d = c.stateNode;
          var e = c.memoizedProps.style;
          e = e !== void 0 && e !== null && e.hasOwnProperty("display") ? e.display : null;
          d.style.display = sb("display", e);
        }
      } else if (c.tag === 6)
        c.stateNode.nodeValue = b ? "" : c.memoizedProps;
      else if ((c.tag !== 23 && c.tag !== 24 || c.memoizedState === null || c === a) && c.child !== null) {
        c.child.return = c;
        c = c.child;
        continue;
      }
      if (c === a)
        break;
      for (; c.sibling === null; ) {
        if (c.return === null || c.return === a)
          return;
        c = c.return;
      }
      c.sibling.return = c.return;
      c = c.sibling;
    }
  }
  function bj(a, b) {
    if (Mf && typeof Mf.onCommitFiberUnmount === "function")
      try {
        Mf.onCommitFiberUnmount(Lf, b);
      } catch (f) {
      }
    switch (b.tag) {
      case 0:
      case 11:
      case 14:
      case 15:
      case 22:
        a = b.updateQueue;
        if (a !== null && (a = a.lastEffect, a !== null)) {
          var c = a = a.next;
          do {
            var d = c, e = d.destroy;
            d = d.tag;
            if (e !== void 0)
              if ((d & 4) !== 0)
                Zi(b, c);
              else {
                d = b;
                try {
                  e();
                } catch (f) {
                  Wi(d, f);
                }
              }
            c = c.next;
          } while (c !== a);
        }
        break;
      case 1:
        Vi(b);
        a = b.stateNode;
        if (typeof a.componentWillUnmount === "function")
          try {
            a.props = b.memoizedProps, a.state = b.memoizedState, a.componentWillUnmount();
          } catch (f) {
            Wi(b, f);
          }
        break;
      case 5:
        Vi(b);
        break;
      case 4:
        cj(a, b);
    }
  }
  function dj(a) {
    a.alternate = null;
    a.child = null;
    a.dependencies = null;
    a.firstEffect = null;
    a.lastEffect = null;
    a.memoizedProps = null;
    a.memoizedState = null;
    a.pendingProps = null;
    a.return = null;
    a.updateQueue = null;
  }
  function ej(a) {
    return a.tag === 5 || a.tag === 3 || a.tag === 4;
  }
  function fj(a) {
    a: {
      for (var b = a.return; b !== null; ) {
        if (ej(b))
          break a;
        b = b.return;
      }
      throw Error(y(160));
    }
    var c = b;
    b = c.stateNode;
    switch (c.tag) {
      case 5:
        var d = false;
        break;
      case 3:
        b = b.containerInfo;
        d = true;
        break;
      case 4:
        b = b.containerInfo;
        d = true;
        break;
      default:
        throw Error(y(161));
    }
    c.flags & 16 && (pb(b, ""), c.flags &= -17);
    a:
      b:
        for (c = a; ; ) {
          for (; c.sibling === null; ) {
            if (c.return === null || ej(c.return)) {
              c = null;
              break a;
            }
            c = c.return;
          }
          c.sibling.return = c.return;
          for (c = c.sibling; c.tag !== 5 && c.tag !== 6 && c.tag !== 18; ) {
            if (c.flags & 2)
              continue b;
            if (c.child === null || c.tag === 4)
              continue b;
            else
              c.child.return = c, c = c.child;
          }
          if (!(c.flags & 2)) {
            c = c.stateNode;
            break a;
          }
        }
    d ? gj(a, c, b) : hj(a, c, b);
  }
  function gj(a, b, c) {
    var d = a.tag, e = d === 5 || d === 6;
    if (e)
      a = e ? a.stateNode : a.stateNode.instance, b ? c.nodeType === 8 ? c.parentNode.insertBefore(a, b) : c.insertBefore(a, b) : (c.nodeType === 8 ? (b = c.parentNode, b.insertBefore(a, c)) : (b = c, b.appendChild(a)), c = c._reactRootContainer, c !== null && c !== void 0 || b.onclick !== null || (b.onclick = jf));
    else if (d !== 4 && (a = a.child, a !== null))
      for (gj(a, b, c), a = a.sibling; a !== null; )
        gj(a, b, c), a = a.sibling;
  }
  function hj(a, b, c) {
    var d = a.tag, e = d === 5 || d === 6;
    if (e)
      a = e ? a.stateNode : a.stateNode.instance, b ? c.insertBefore(a, b) : c.appendChild(a);
    else if (d !== 4 && (a = a.child, a !== null))
      for (hj(a, b, c), a = a.sibling; a !== null; )
        hj(a, b, c), a = a.sibling;
  }
  function cj(a, b) {
    for (var c = b, d = false, e, f; ; ) {
      if (!d) {
        d = c.return;
        a:
          for (; ; ) {
            if (d === null)
              throw Error(y(160));
            e = d.stateNode;
            switch (d.tag) {
              case 5:
                f = false;
                break a;
              case 3:
                e = e.containerInfo;
                f = true;
                break a;
              case 4:
                e = e.containerInfo;
                f = true;
                break a;
            }
            d = d.return;
          }
        d = true;
      }
      if (c.tag === 5 || c.tag === 6) {
        a:
          for (var g = a, h = c, k = h; ; )
            if (bj(g, k), k.child !== null && k.tag !== 4)
              k.child.return = k, k = k.child;
            else {
              if (k === h)
                break a;
              for (; k.sibling === null; ) {
                if (k.return === null || k.return === h)
                  break a;
                k = k.return;
              }
              k.sibling.return = k.return;
              k = k.sibling;
            }
        f ? (g = e, h = c.stateNode, g.nodeType === 8 ? g.parentNode.removeChild(h) : g.removeChild(h)) : e.removeChild(c.stateNode);
      } else if (c.tag === 4) {
        if (c.child !== null) {
          e = c.stateNode.containerInfo;
          f = true;
          c.child.return = c;
          c = c.child;
          continue;
        }
      } else if (bj(a, c), c.child !== null) {
        c.child.return = c;
        c = c.child;
        continue;
      }
      if (c === b)
        break;
      for (; c.sibling === null; ) {
        if (c.return === null || c.return === b)
          return;
        c = c.return;
        c.tag === 4 && (d = false);
      }
      c.sibling.return = c.return;
      c = c.sibling;
    }
  }
  function ij(a, b) {
    switch (b.tag) {
      case 0:
      case 11:
      case 14:
      case 15:
      case 22:
        var c = b.updateQueue;
        c = c !== null ? c.lastEffect : null;
        if (c !== null) {
          var d = c = c.next;
          do
            (d.tag & 3) === 3 && (a = d.destroy, d.destroy = void 0, a !== void 0 && a()), d = d.next;
          while (d !== c);
        }
        return;
      case 1:
        return;
      case 5:
        c = b.stateNode;
        if (c != null) {
          d = b.memoizedProps;
          var e = a !== null ? a.memoizedProps : d;
          a = b.type;
          var f = b.updateQueue;
          b.updateQueue = null;
          if (f !== null) {
            c[xf] = d;
            a === "input" && d.type === "radio" && d.name != null && $a(c, d);
            wb(a, e);
            b = wb(a, d);
            for (e = 0; e < f.length; e += 2) {
              var g = f[e], h = f[e + 1];
              g === "style" ? tb(c, h) : g === "dangerouslySetInnerHTML" ? ob(c, h) : g === "children" ? pb(c, h) : qa(c, g, h, b);
            }
            switch (a) {
              case "input":
                ab(c, d);
                break;
              case "textarea":
                ib(c, d);
                break;
              case "select":
                a = c._wrapperState.wasMultiple, c._wrapperState.wasMultiple = !!d.multiple, f = d.value, f != null ? fb(c, !!d.multiple, f, false) : a !== !!d.multiple && (d.defaultValue != null ? fb(c, !!d.multiple, d.defaultValue, true) : fb(c, !!d.multiple, d.multiple ? [] : "", false));
            }
          }
        }
        return;
      case 6:
        if (b.stateNode === null)
          throw Error(y(162));
        b.stateNode.nodeValue = b.memoizedProps;
        return;
      case 3:
        c = b.stateNode;
        c.hydrate && (c.hydrate = false, Cc(c.containerInfo));
        return;
      case 12:
        return;
      case 13:
        b.memoizedState !== null && (jj = O(), aj(b.child, true));
        kj(b);
        return;
      case 19:
        kj(b);
        return;
      case 17:
        return;
      case 23:
      case 24:
        aj(b, b.memoizedState !== null);
        return;
    }
    throw Error(y(163));
  }
  function kj(a) {
    var b = a.updateQueue;
    if (b !== null) {
      a.updateQueue = null;
      var c = a.stateNode;
      c === null && (c = a.stateNode = new Ui());
      b.forEach(function(b2) {
        var d = lj.bind(null, a, b2);
        c.has(b2) || (c.add(b2), b2.then(d, d));
      });
    }
  }
  function mj(a, b) {
    return a !== null && (a = a.memoizedState, a === null || a.dehydrated !== null) ? (b = b.memoizedState, b !== null && b.dehydrated === null) : false;
  }
  var nj = Math.ceil;
  var oj = ra.ReactCurrentDispatcher;
  var pj = ra.ReactCurrentOwner;
  var X = 0;
  var U = null;
  var Y = null;
  var W = 0;
  var qj = 0;
  var rj = Bf(0);
  var V = 0;
  var sj = null;
  var tj = 0;
  var Dg = 0;
  var Hi = 0;
  var uj = 0;
  var vj = null;
  var jj = 0;
  var Ji = Infinity;
  function wj() {
    Ji = O() + 500;
  }
  var Z = null;
  var Qi = false;
  var Ri = null;
  var Ti = null;
  var xj = false;
  var yj = null;
  var zj = 90;
  var Aj = [];
  var Bj = [];
  var Cj = null;
  var Dj = 0;
  var Ej = null;
  var Fj = -1;
  var Gj = 0;
  var Hj = 0;
  var Ij = null;
  var Jj = false;
  function Hg() {
    return (X & 48) !== 0 ? O() : Fj !== -1 ? Fj : Fj = O();
  }
  function Ig(a) {
    a = a.mode;
    if ((a & 2) === 0)
      return 1;
    if ((a & 4) === 0)
      return eg() === 99 ? 1 : 2;
    Gj === 0 && (Gj = tj);
    if (kg.transition !== 0) {
      Hj !== 0 && (Hj = vj !== null ? vj.pendingLanes : 0);
      a = Gj;
      var b = 4186112 & ~Hj;
      b &= -b;
      b === 0 && (a = 4186112 & ~a, b = a & -a, b === 0 && (b = 8192));
      return b;
    }
    a = eg();
    (X & 4) !== 0 && a === 98 ? a = Xc(12, Gj) : (a = Sc(a), a = Xc(a, Gj));
    return a;
  }
  function Jg(a, b, c) {
    if (50 < Dj)
      throw Dj = 0, Ej = null, Error(y(185));
    a = Kj(a, b);
    if (a === null)
      return null;
    $c(a, b, c);
    a === U && (Hi |= b, V === 4 && Ii(a, W));
    var d = eg();
    b === 1 ? (X & 8) !== 0 && (X & 48) === 0 ? Lj(a) : (Mj(a, c), X === 0 && (wj(), ig())) : ((X & 4) === 0 || d !== 98 && d !== 99 || (Cj === null ? Cj = new Set([a]) : Cj.add(a)), Mj(a, c));
    vj = a;
  }
  function Kj(a, b) {
    a.lanes |= b;
    var c = a.alternate;
    c !== null && (c.lanes |= b);
    c = a;
    for (a = a.return; a !== null; )
      a.childLanes |= b, c = a.alternate, c !== null && (c.childLanes |= b), c = a, a = a.return;
    return c.tag === 3 ? c.stateNode : null;
  }
  function Mj(a, b) {
    for (var c = a.callbackNode, d = a.suspendedLanes, e = a.pingedLanes, f = a.expirationTimes, g = a.pendingLanes; 0 < g; ) {
      var h = 31 - Vc(g), k = 1 << h, l = f[h];
      if (l === -1) {
        if ((k & d) === 0 || (k & e) !== 0) {
          l = b;
          Rc(k);
          var n = F;
          f[h] = 10 <= n ? l + 250 : 6 <= n ? l + 5e3 : -1;
        }
      } else
        l <= b && (a.expiredLanes |= k);
      g &= ~k;
    }
    d = Uc(a, a === U ? W : 0);
    b = F;
    if (d === 0)
      c !== null && (c !== Zf && Pf(c), a.callbackNode = null, a.callbackPriority = 0);
    else {
      if (c !== null) {
        if (a.callbackPriority === b)
          return;
        c !== Zf && Pf(c);
      }
      b === 15 ? (c = Lj.bind(null, a), ag === null ? (ag = [c], bg = Of(Uf, jg)) : ag.push(c), c = Zf) : b === 14 ? c = hg(99, Lj.bind(null, a)) : (c = Tc(b), c = hg(c, Nj.bind(null, a)));
      a.callbackPriority = b;
      a.callbackNode = c;
    }
  }
  function Nj(a) {
    Fj = -1;
    Hj = Gj = 0;
    if ((X & 48) !== 0)
      throw Error(y(327));
    var b = a.callbackNode;
    if (Oj() && a.callbackNode !== b)
      return null;
    var c = Uc(a, a === U ? W : 0);
    if (c === 0)
      return null;
    var d = c;
    var e = X;
    X |= 16;
    var f = Pj();
    if (U !== a || W !== d)
      wj(), Qj(a, d);
    do
      try {
        Rj();
        break;
      } catch (h) {
        Sj(a, h);
      }
    while (1);
    qg();
    oj.current = f;
    X = e;
    Y !== null ? d = 0 : (U = null, W = 0, d = V);
    if ((tj & Hi) !== 0)
      Qj(a, 0);
    else if (d !== 0) {
      d === 2 && (X |= 64, a.hydrate && (a.hydrate = false, qf(a.containerInfo)), c = Wc(a), c !== 0 && (d = Tj(a, c)));
      if (d === 1)
        throw b = sj, Qj(a, 0), Ii(a, c), Mj(a, O()), b;
      a.finishedWork = a.current.alternate;
      a.finishedLanes = c;
      switch (d) {
        case 0:
        case 1:
          throw Error(y(345));
        case 2:
          Uj(a);
          break;
        case 3:
          Ii(a, c);
          if ((c & 62914560) === c && (d = jj + 500 - O(), 10 < d)) {
            if (Uc(a, 0) !== 0)
              break;
            e = a.suspendedLanes;
            if ((e & c) !== c) {
              Hg();
              a.pingedLanes |= a.suspendedLanes & e;
              break;
            }
            a.timeoutHandle = of(Uj.bind(null, a), d);
            break;
          }
          Uj(a);
          break;
        case 4:
          Ii(a, c);
          if ((c & 4186112) === c)
            break;
          d = a.eventTimes;
          for (e = -1; 0 < c; ) {
            var g = 31 - Vc(c);
            f = 1 << g;
            g = d[g];
            g > e && (e = g);
            c &= ~f;
          }
          c = e;
          c = O() - c;
          c = (120 > c ? 120 : 480 > c ? 480 : 1080 > c ? 1080 : 1920 > c ? 1920 : 3e3 > c ? 3e3 : 4320 > c ? 4320 : 1960 * nj(c / 1960)) - c;
          if (10 < c) {
            a.timeoutHandle = of(Uj.bind(null, a), c);
            break;
          }
          Uj(a);
          break;
        case 5:
          Uj(a);
          break;
        default:
          throw Error(y(329));
      }
    }
    Mj(a, O());
    return a.callbackNode === b ? Nj.bind(null, a) : null;
  }
  function Ii(a, b) {
    b &= ~uj;
    b &= ~Hi;
    a.suspendedLanes |= b;
    a.pingedLanes &= ~b;
    for (a = a.expirationTimes; 0 < b; ) {
      var c = 31 - Vc(b), d = 1 << c;
      a[c] = -1;
      b &= ~d;
    }
  }
  function Lj(a) {
    if ((X & 48) !== 0)
      throw Error(y(327));
    Oj();
    if (a === U && (a.expiredLanes & W) !== 0) {
      var b = W;
      var c = Tj(a, b);
      (tj & Hi) !== 0 && (b = Uc(a, b), c = Tj(a, b));
    } else
      b = Uc(a, 0), c = Tj(a, b);
    a.tag !== 0 && c === 2 && (X |= 64, a.hydrate && (a.hydrate = false, qf(a.containerInfo)), b = Wc(a), b !== 0 && (c = Tj(a, b)));
    if (c === 1)
      throw c = sj, Qj(a, 0), Ii(a, b), Mj(a, O()), c;
    a.finishedWork = a.current.alternate;
    a.finishedLanes = b;
    Uj(a);
    Mj(a, O());
    return null;
  }
  function Vj() {
    if (Cj !== null) {
      var a = Cj;
      Cj = null;
      a.forEach(function(a2) {
        a2.expiredLanes |= 24 & a2.pendingLanes;
        Mj(a2, O());
      });
    }
    ig();
  }
  function Wj(a, b) {
    var c = X;
    X |= 1;
    try {
      return a(b);
    } finally {
      X = c, X === 0 && (wj(), ig());
    }
  }
  function Xj(a, b) {
    var c = X;
    X &= -2;
    X |= 8;
    try {
      return a(b);
    } finally {
      X = c, X === 0 && (wj(), ig());
    }
  }
  function ni(a, b) {
    I(rj, qj);
    qj |= b;
    tj |= b;
  }
  function Ki() {
    qj = rj.current;
    H(rj);
  }
  function Qj(a, b) {
    a.finishedWork = null;
    a.finishedLanes = 0;
    var c = a.timeoutHandle;
    c !== -1 && (a.timeoutHandle = -1, pf(c));
    if (Y !== null)
      for (c = Y.return; c !== null; ) {
        var d = c;
        switch (d.tag) {
          case 1:
            d = d.type.childContextTypes;
            d !== null && d !== void 0 && Gf();
            break;
          case 3:
            fh();
            H(N);
            H(M);
            uh();
            break;
          case 5:
            hh(d);
            break;
          case 4:
            fh();
            break;
          case 13:
            H(P);
            break;
          case 19:
            H(P);
            break;
          case 10:
            rg(d);
            break;
          case 23:
          case 24:
            Ki();
        }
        c = c.return;
      }
    U = a;
    Y = Tg(a.current, null);
    W = qj = tj = b;
    V = 0;
    sj = null;
    uj = Hi = Dg = 0;
  }
  function Sj(a, b) {
    do {
      var c = Y;
      try {
        qg();
        vh.current = Gh;
        if (yh) {
          for (var d = R.memoizedState; d !== null; ) {
            var e = d.queue;
            e !== null && (e.pending = null);
            d = d.next;
          }
          yh = false;
        }
        xh = 0;
        T = S = R = null;
        zh = false;
        pj.current = null;
        if (c === null || c.return === null) {
          V = 1;
          sj = b;
          Y = null;
          break;
        }
        a: {
          var f = a, g = c.return, h = c, k = b;
          b = W;
          h.flags |= 2048;
          h.firstEffect = h.lastEffect = null;
          if (k !== null && typeof k === "object" && typeof k.then === "function") {
            var l = k;
            if ((h.mode & 2) === 0) {
              var n = h.alternate;
              n ? (h.updateQueue = n.updateQueue, h.memoizedState = n.memoizedState, h.lanes = n.lanes) : (h.updateQueue = null, h.memoizedState = null);
            }
            var A = (P.current & 1) !== 0, p = g;
            do {
              var C;
              if (C = p.tag === 13) {
                var x = p.memoizedState;
                if (x !== null)
                  C = x.dehydrated !== null ? true : false;
                else {
                  var w = p.memoizedProps;
                  C = w.fallback === void 0 ? false : w.unstable_avoidThisFallback !== true ? true : A ? false : true;
                }
              }
              if (C) {
                var z = p.updateQueue;
                if (z === null) {
                  var u = new Set();
                  u.add(l);
                  p.updateQueue = u;
                } else
                  z.add(l);
                if ((p.mode & 2) === 0) {
                  p.flags |= 64;
                  h.flags |= 16384;
                  h.flags &= -2981;
                  if (h.tag === 1)
                    if (h.alternate === null)
                      h.tag = 17;
                    else {
                      var t = zg(-1, 1);
                      t.tag = 2;
                      Ag(h, t);
                    }
                  h.lanes |= 1;
                  break a;
                }
                k = void 0;
                h = b;
                var q = f.pingCache;
                q === null ? (q = f.pingCache = new Oi(), k = new Set(), q.set(l, k)) : (k = q.get(l), k === void 0 && (k = new Set(), q.set(l, k)));
                if (!k.has(h)) {
                  k.add(h);
                  var v = Yj.bind(null, f, l, h);
                  l.then(v, v);
                }
                p.flags |= 4096;
                p.lanes = b;
                break a;
              }
              p = p.return;
            } while (p !== null);
            k = Error((Ra(h.type) || "A React component") + " suspended while rendering, but no fallback UI was specified.\n\nAdd a <Suspense fallback=...> component higher in the tree to provide a loading indicator or placeholder to display.");
          }
          V !== 5 && (V = 2);
          k = Mi(k, h);
          p = g;
          do {
            switch (p.tag) {
              case 3:
                f = k;
                p.flags |= 4096;
                b &= -b;
                p.lanes |= b;
                var J = Pi(p, f, b);
                Bg(p, J);
                break a;
              case 1:
                f = k;
                var K = p.type, Q = p.stateNode;
                if ((p.flags & 64) === 0 && (typeof K.getDerivedStateFromError === "function" || Q !== null && typeof Q.componentDidCatch === "function" && (Ti === null || !Ti.has(Q)))) {
                  p.flags |= 4096;
                  b &= -b;
                  p.lanes |= b;
                  var L = Si(p, f, b);
                  Bg(p, L);
                  break a;
                }
            }
            p = p.return;
          } while (p !== null);
        }
        Zj(c);
      } catch (va) {
        b = va;
        Y === c && c !== null && (Y = c = c.return);
        continue;
      }
      break;
    } while (1);
  }
  function Pj() {
    var a = oj.current;
    oj.current = Gh;
    return a === null ? Gh : a;
  }
  function Tj(a, b) {
    var c = X;
    X |= 16;
    var d = Pj();
    U === a && W === b || Qj(a, b);
    do
      try {
        ak();
        break;
      } catch (e) {
        Sj(a, e);
      }
    while (1);
    qg();
    X = c;
    oj.current = d;
    if (Y !== null)
      throw Error(y(261));
    U = null;
    W = 0;
    return V;
  }
  function ak() {
    for (; Y !== null; )
      bk(Y);
  }
  function Rj() {
    for (; Y !== null && !Qf(); )
      bk(Y);
  }
  function bk(a) {
    var b = ck(a.alternate, a, qj);
    a.memoizedProps = a.pendingProps;
    b === null ? Zj(a) : Y = b;
    pj.current = null;
  }
  function Zj(a) {
    var b = a;
    do {
      var c = b.alternate;
      a = b.return;
      if ((b.flags & 2048) === 0) {
        c = Gi(c, b, qj);
        if (c !== null) {
          Y = c;
          return;
        }
        c = b;
        if (c.tag !== 24 && c.tag !== 23 || c.memoizedState === null || (qj & 1073741824) !== 0 || (c.mode & 4) === 0) {
          for (var d = 0, e = c.child; e !== null; )
            d |= e.lanes | e.childLanes, e = e.sibling;
          c.childLanes = d;
        }
        a !== null && (a.flags & 2048) === 0 && (a.firstEffect === null && (a.firstEffect = b.firstEffect), b.lastEffect !== null && (a.lastEffect !== null && (a.lastEffect.nextEffect = b.firstEffect), a.lastEffect = b.lastEffect), 1 < b.flags && (a.lastEffect !== null ? a.lastEffect.nextEffect = b : a.firstEffect = b, a.lastEffect = b));
      } else {
        c = Li(b);
        if (c !== null) {
          c.flags &= 2047;
          Y = c;
          return;
        }
        a !== null && (a.firstEffect = a.lastEffect = null, a.flags |= 2048);
      }
      b = b.sibling;
      if (b !== null) {
        Y = b;
        return;
      }
      Y = b = a;
    } while (b !== null);
    V === 0 && (V = 5);
  }
  function Uj(a) {
    var b = eg();
    gg(99, dk.bind(null, a, b));
    return null;
  }
  function dk(a, b) {
    do
      Oj();
    while (yj !== null);
    if ((X & 48) !== 0)
      throw Error(y(327));
    var c = a.finishedWork;
    if (c === null)
      return null;
    a.finishedWork = null;
    a.finishedLanes = 0;
    if (c === a.current)
      throw Error(y(177));
    a.callbackNode = null;
    var d = c.lanes | c.childLanes, e = d, f = a.pendingLanes & ~e;
    a.pendingLanes = e;
    a.suspendedLanes = 0;
    a.pingedLanes = 0;
    a.expiredLanes &= e;
    a.mutableReadLanes &= e;
    a.entangledLanes &= e;
    e = a.entanglements;
    for (var g = a.eventTimes, h = a.expirationTimes; 0 < f; ) {
      var k = 31 - Vc(f), l = 1 << k;
      e[k] = 0;
      g[k] = -1;
      h[k] = -1;
      f &= ~l;
    }
    Cj !== null && (d & 24) === 0 && Cj.has(a) && Cj.delete(a);
    a === U && (Y = U = null, W = 0);
    1 < c.flags ? c.lastEffect !== null ? (c.lastEffect.nextEffect = c, d = c.firstEffect) : d = c : d = c.firstEffect;
    if (d !== null) {
      e = X;
      X |= 32;
      pj.current = null;
      kf = fd;
      g = Ne();
      if (Oe(g)) {
        if ("selectionStart" in g)
          h = {start: g.selectionStart, end: g.selectionEnd};
        else
          a:
            if (h = (h = g.ownerDocument) && h.defaultView || window, (l = h.getSelection && h.getSelection()) && l.rangeCount !== 0) {
              h = l.anchorNode;
              f = l.anchorOffset;
              k = l.focusNode;
              l = l.focusOffset;
              try {
                h.nodeType, k.nodeType;
              } catch (va) {
                h = null;
                break a;
              }
              var n = 0, A = -1, p = -1, C = 0, x = 0, w = g, z = null;
              b:
                for (; ; ) {
                  for (var u; ; ) {
                    w !== h || f !== 0 && w.nodeType !== 3 || (A = n + f);
                    w !== k || l !== 0 && w.nodeType !== 3 || (p = n + l);
                    w.nodeType === 3 && (n += w.nodeValue.length);
                    if ((u = w.firstChild) === null)
                      break;
                    z = w;
                    w = u;
                  }
                  for (; ; ) {
                    if (w === g)
                      break b;
                    z === h && ++C === f && (A = n);
                    z === k && ++x === l && (p = n);
                    if ((u = w.nextSibling) !== null)
                      break;
                    w = z;
                    z = w.parentNode;
                  }
                  w = u;
                }
              h = A === -1 || p === -1 ? null : {start: A, end: p};
            } else
              h = null;
        h = h || {start: 0, end: 0};
      } else
        h = null;
      lf = {focusedElem: g, selectionRange: h};
      fd = false;
      Ij = null;
      Jj = false;
      Z = d;
      do
        try {
          ek();
        } catch (va) {
          if (Z === null)
            throw Error(y(330));
          Wi(Z, va);
          Z = Z.nextEffect;
        }
      while (Z !== null);
      Ij = null;
      Z = d;
      do
        try {
          for (g = a; Z !== null; ) {
            var t = Z.flags;
            t & 16 && pb(Z.stateNode, "");
            if (t & 128) {
              var q = Z.alternate;
              if (q !== null) {
                var v = q.ref;
                v !== null && (typeof v === "function" ? v(null) : v.current = null);
              }
            }
            switch (t & 1038) {
              case 2:
                fj(Z);
                Z.flags &= -3;
                break;
              case 6:
                fj(Z);
                Z.flags &= -3;
                ij(Z.alternate, Z);
                break;
              case 1024:
                Z.flags &= -1025;
                break;
              case 1028:
                Z.flags &= -1025;
                ij(Z.alternate, Z);
                break;
              case 4:
                ij(Z.alternate, Z);
                break;
              case 8:
                h = Z;
                cj(g, h);
                var J = h.alternate;
                dj(h);
                J !== null && dj(J);
            }
            Z = Z.nextEffect;
          }
        } catch (va) {
          if (Z === null)
            throw Error(y(330));
          Wi(Z, va);
          Z = Z.nextEffect;
        }
      while (Z !== null);
      v = lf;
      q = Ne();
      t = v.focusedElem;
      g = v.selectionRange;
      if (q !== t && t && t.ownerDocument && Me(t.ownerDocument.documentElement, t)) {
        g !== null && Oe(t) && (q = g.start, v = g.end, v === void 0 && (v = q), "selectionStart" in t ? (t.selectionStart = q, t.selectionEnd = Math.min(v, t.value.length)) : (v = (q = t.ownerDocument || document) && q.defaultView || window, v.getSelection && (v = v.getSelection(), h = t.textContent.length, J = Math.min(g.start, h), g = g.end === void 0 ? J : Math.min(g.end, h), !v.extend && J > g && (h = g, g = J, J = h), h = Le(t, J), f = Le(t, g), h && f && (v.rangeCount !== 1 || v.anchorNode !== h.node || v.anchorOffset !== h.offset || v.focusNode !== f.node || v.focusOffset !== f.offset) && (q = q.createRange(), q.setStart(h.node, h.offset), v.removeAllRanges(), J > g ? (v.addRange(q), v.extend(f.node, f.offset)) : (q.setEnd(f.node, f.offset), v.addRange(q))))));
        q = [];
        for (v = t; v = v.parentNode; )
          v.nodeType === 1 && q.push({element: v, left: v.scrollLeft, top: v.scrollTop});
        typeof t.focus === "function" && t.focus();
        for (t = 0; t < q.length; t++)
          v = q[t], v.element.scrollLeft = v.left, v.element.scrollTop = v.top;
      }
      fd = !!kf;
      lf = kf = null;
      a.current = c;
      Z = d;
      do
        try {
          for (t = a; Z !== null; ) {
            var K = Z.flags;
            K & 36 && Yi(t, Z.alternate, Z);
            if (K & 128) {
              q = void 0;
              var Q = Z.ref;
              if (Q !== null) {
                var L = Z.stateNode;
                switch (Z.tag) {
                  case 5:
                    q = L;
                    break;
                  default:
                    q = L;
                }
                typeof Q === "function" ? Q(q) : Q.current = q;
              }
            }
            Z = Z.nextEffect;
          }
        } catch (va) {
          if (Z === null)
            throw Error(y(330));
          Wi(Z, va);
          Z = Z.nextEffect;
        }
      while (Z !== null);
      Z = null;
      $f();
      X = e;
    } else
      a.current = c;
    if (xj)
      xj = false, yj = a, zj = b;
    else
      for (Z = d; Z !== null; )
        b = Z.nextEffect, Z.nextEffect = null, Z.flags & 8 && (K = Z, K.sibling = null, K.stateNode = null), Z = b;
    d = a.pendingLanes;
    d === 0 && (Ti = null);
    d === 1 ? a === Ej ? Dj++ : (Dj = 0, Ej = a) : Dj = 0;
    c = c.stateNode;
    if (Mf && typeof Mf.onCommitFiberRoot === "function")
      try {
        Mf.onCommitFiberRoot(Lf, c, void 0, (c.current.flags & 64) === 64);
      } catch (va) {
      }
    Mj(a, O());
    if (Qi)
      throw Qi = false, a = Ri, Ri = null, a;
    if ((X & 8) !== 0)
      return null;
    ig();
    return null;
  }
  function ek() {
    for (; Z !== null; ) {
      var a = Z.alternate;
      Jj || Ij === null || ((Z.flags & 8) !== 0 ? dc(Z, Ij) && (Jj = true) : Z.tag === 13 && mj(a, Z) && dc(Z, Ij) && (Jj = true));
      var b = Z.flags;
      (b & 256) !== 0 && Xi(a, Z);
      (b & 512) === 0 || xj || (xj = true, hg(97, function() {
        Oj();
        return null;
      }));
      Z = Z.nextEffect;
    }
  }
  function Oj() {
    if (zj !== 90) {
      var a = 97 < zj ? 97 : zj;
      zj = 90;
      return gg(a, fk);
    }
    return false;
  }
  function $i(a, b) {
    Aj.push(b, a);
    xj || (xj = true, hg(97, function() {
      Oj();
      return null;
    }));
  }
  function Zi(a, b) {
    Bj.push(b, a);
    xj || (xj = true, hg(97, function() {
      Oj();
      return null;
    }));
  }
  function fk() {
    if (yj === null)
      return false;
    var a = yj;
    yj = null;
    if ((X & 48) !== 0)
      throw Error(y(331));
    var b = X;
    X |= 32;
    var c = Bj;
    Bj = [];
    for (var d = 0; d < c.length; d += 2) {
      var e = c[d], f = c[d + 1], g = e.destroy;
      e.destroy = void 0;
      if (typeof g === "function")
        try {
          g();
        } catch (k) {
          if (f === null)
            throw Error(y(330));
          Wi(f, k);
        }
    }
    c = Aj;
    Aj = [];
    for (d = 0; d < c.length; d += 2) {
      e = c[d];
      f = c[d + 1];
      try {
        var h = e.create;
        e.destroy = h();
      } catch (k) {
        if (f === null)
          throw Error(y(330));
        Wi(f, k);
      }
    }
    for (h = a.current.firstEffect; h !== null; )
      a = h.nextEffect, h.nextEffect = null, h.flags & 8 && (h.sibling = null, h.stateNode = null), h = a;
    X = b;
    ig();
    return true;
  }
  function gk(a, b, c) {
    b = Mi(c, b);
    b = Pi(a, b, 1);
    Ag(a, b);
    b = Hg();
    a = Kj(a, 1);
    a !== null && ($c(a, 1, b), Mj(a, b));
  }
  function Wi(a, b) {
    if (a.tag === 3)
      gk(a, a, b);
    else
      for (var c = a.return; c !== null; ) {
        if (c.tag === 3) {
          gk(c, a, b);
          break;
        } else if (c.tag === 1) {
          var d = c.stateNode;
          if (typeof c.type.getDerivedStateFromError === "function" || typeof d.componentDidCatch === "function" && (Ti === null || !Ti.has(d))) {
            a = Mi(b, a);
            var e = Si(c, a, 1);
            Ag(c, e);
            e = Hg();
            c = Kj(c, 1);
            if (c !== null)
              $c(c, 1, e), Mj(c, e);
            else if (typeof d.componentDidCatch === "function" && (Ti === null || !Ti.has(d)))
              try {
                d.componentDidCatch(b, a);
              } catch (f) {
              }
            break;
          }
        }
        c = c.return;
      }
  }
  function Yj(a, b, c) {
    var d = a.pingCache;
    d !== null && d.delete(b);
    b = Hg();
    a.pingedLanes |= a.suspendedLanes & c;
    U === a && (W & c) === c && (V === 4 || V === 3 && (W & 62914560) === W && 500 > O() - jj ? Qj(a, 0) : uj |= c);
    Mj(a, b);
  }
  function lj(a, b) {
    var c = a.stateNode;
    c !== null && c.delete(b);
    b = 0;
    b === 0 && (b = a.mode, (b & 2) === 0 ? b = 1 : (b & 4) === 0 ? b = eg() === 99 ? 1 : 2 : (Gj === 0 && (Gj = tj), b = Yc(62914560 & ~Gj), b === 0 && (b = 4194304)));
    c = Hg();
    a = Kj(a, b);
    a !== null && ($c(a, b, c), Mj(a, c));
  }
  var ck;
  ck = function(a, b, c) {
    var d = b.lanes;
    if (a !== null)
      if (a.memoizedProps !== b.pendingProps || N.current)
        ug = true;
      else if ((c & d) !== 0)
        ug = (a.flags & 16384) !== 0 ? true : false;
      else {
        ug = false;
        switch (b.tag) {
          case 3:
            ri(b);
            sh();
            break;
          case 5:
            gh(b);
            break;
          case 1:
            Ff(b.type) && Jf(b);
            break;
          case 4:
            eh(b, b.stateNode.containerInfo);
            break;
          case 10:
            d = b.memoizedProps.value;
            var e = b.type._context;
            I(mg, e._currentValue);
            e._currentValue = d;
            break;
          case 13:
            if (b.memoizedState !== null) {
              if ((c & b.child.childLanes) !== 0)
                return ti(a, b, c);
              I(P, P.current & 1);
              b = hi(a, b, c);
              return b !== null ? b.sibling : null;
            }
            I(P, P.current & 1);
            break;
          case 19:
            d = (c & b.childLanes) !== 0;
            if ((a.flags & 64) !== 0) {
              if (d)
                return Ai(a, b, c);
              b.flags |= 64;
            }
            e = b.memoizedState;
            e !== null && (e.rendering = null, e.tail = null, e.lastEffect = null);
            I(P, P.current);
            if (d)
              break;
            else
              return null;
          case 23:
          case 24:
            return b.lanes = 0, mi(a, b, c);
        }
        return hi(a, b, c);
      }
    else
      ug = false;
    b.lanes = 0;
    switch (b.tag) {
      case 2:
        d = b.type;
        a !== null && (a.alternate = null, b.alternate = null, b.flags |= 2);
        a = b.pendingProps;
        e = Ef(b, M.current);
        tg(b, c);
        e = Ch(null, b, d, a, e, c);
        b.flags |= 1;
        if (typeof e === "object" && e !== null && typeof e.render === "function" && e.$$typeof === void 0) {
          b.tag = 1;
          b.memoizedState = null;
          b.updateQueue = null;
          if (Ff(d)) {
            var f = true;
            Jf(b);
          } else
            f = false;
          b.memoizedState = e.state !== null && e.state !== void 0 ? e.state : null;
          xg(b);
          var g = d.getDerivedStateFromProps;
          typeof g === "function" && Gg(b, d, g, a);
          e.updater = Kg;
          b.stateNode = e;
          e._reactInternals = b;
          Og(b, d, a, c);
          b = qi(null, b, d, true, f, c);
        } else
          b.tag = 0, fi(null, b, e, c), b = b.child;
        return b;
      case 16:
        e = b.elementType;
        a: {
          a !== null && (a.alternate = null, b.alternate = null, b.flags |= 2);
          a = b.pendingProps;
          f = e._init;
          e = f(e._payload);
          b.type = e;
          f = b.tag = hk(e);
          a = lg(e, a);
          switch (f) {
            case 0:
              b = li(null, b, e, a, c);
              break a;
            case 1:
              b = pi(null, b, e, a, c);
              break a;
            case 11:
              b = gi(null, b, e, a, c);
              break a;
            case 14:
              b = ii(null, b, e, lg(e.type, a), d, c);
              break a;
          }
          throw Error(y(306, e, ""));
        }
        return b;
      case 0:
        return d = b.type, e = b.pendingProps, e = b.elementType === d ? e : lg(d, e), li(a, b, d, e, c);
      case 1:
        return d = b.type, e = b.pendingProps, e = b.elementType === d ? e : lg(d, e), pi(a, b, d, e, c);
      case 3:
        ri(b);
        d = b.updateQueue;
        if (a === null || d === null)
          throw Error(y(282));
        d = b.pendingProps;
        e = b.memoizedState;
        e = e !== null ? e.element : null;
        yg(a, b);
        Cg(b, d, null, c);
        d = b.memoizedState.element;
        if (d === e)
          sh(), b = hi(a, b, c);
        else {
          e = b.stateNode;
          if (f = e.hydrate)
            kh = rf(b.stateNode.containerInfo.firstChild), jh = b, f = lh = true;
          if (f) {
            a = e.mutableSourceEagerHydrationData;
            if (a != null)
              for (e = 0; e < a.length; e += 2)
                f = a[e], f._workInProgressVersionPrimary = a[e + 1], th.push(f);
            c = Zg(b, null, d, c);
            for (b.child = c; c; )
              c.flags = c.flags & -3 | 1024, c = c.sibling;
          } else
            fi(a, b, d, c), sh();
          b = b.child;
        }
        return b;
      case 5:
        return gh(b), a === null && ph(b), d = b.type, e = b.pendingProps, f = a !== null ? a.memoizedProps : null, g = e.children, nf(d, e) ? g = null : f !== null && nf(d, f) && (b.flags |= 16), oi(a, b), fi(a, b, g, c), b.child;
      case 6:
        return a === null && ph(b), null;
      case 13:
        return ti(a, b, c);
      case 4:
        return eh(b, b.stateNode.containerInfo), d = b.pendingProps, a === null ? b.child = Yg(b, null, d, c) : fi(a, b, d, c), b.child;
      case 11:
        return d = b.type, e = b.pendingProps, e = b.elementType === d ? e : lg(d, e), gi(a, b, d, e, c);
      case 7:
        return fi(a, b, b.pendingProps, c), b.child;
      case 8:
        return fi(a, b, b.pendingProps.children, c), b.child;
      case 12:
        return fi(a, b, b.pendingProps.children, c), b.child;
      case 10:
        a: {
          d = b.type._context;
          e = b.pendingProps;
          g = b.memoizedProps;
          f = e.value;
          var h = b.type._context;
          I(mg, h._currentValue);
          h._currentValue = f;
          if (g !== null)
            if (h = g.value, f = He(h, f) ? 0 : (typeof d._calculateChangedBits === "function" ? d._calculateChangedBits(h, f) : 1073741823) | 0, f === 0) {
              if (g.children === e.children && !N.current) {
                b = hi(a, b, c);
                break a;
              }
            } else
              for (h = b.child, h !== null && (h.return = b); h !== null; ) {
                var k = h.dependencies;
                if (k !== null) {
                  g = h.child;
                  for (var l = k.firstContext; l !== null; ) {
                    if (l.context === d && (l.observedBits & f) !== 0) {
                      h.tag === 1 && (l = zg(-1, c & -c), l.tag = 2, Ag(h, l));
                      h.lanes |= c;
                      l = h.alternate;
                      l !== null && (l.lanes |= c);
                      sg(h.return, c);
                      k.lanes |= c;
                      break;
                    }
                    l = l.next;
                  }
                } else
                  g = h.tag === 10 ? h.type === b.type ? null : h.child : h.child;
                if (g !== null)
                  g.return = h;
                else
                  for (g = h; g !== null; ) {
                    if (g === b) {
                      g = null;
                      break;
                    }
                    h = g.sibling;
                    if (h !== null) {
                      h.return = g.return;
                      g = h;
                      break;
                    }
                    g = g.return;
                  }
                h = g;
              }
          fi(a, b, e.children, c);
          b = b.child;
        }
        return b;
      case 9:
        return e = b.type, f = b.pendingProps, d = f.children, tg(b, c), e = vg(e, f.unstable_observedBits), d = d(e), b.flags |= 1, fi(a, b, d, c), b.child;
      case 14:
        return e = b.type, f = lg(e, b.pendingProps), f = lg(e.type, f), ii(a, b, e, f, d, c);
      case 15:
        return ki(a, b, b.type, b.pendingProps, d, c);
      case 17:
        return d = b.type, e = b.pendingProps, e = b.elementType === d ? e : lg(d, e), a !== null && (a.alternate = null, b.alternate = null, b.flags |= 2), b.tag = 1, Ff(d) ? (a = true, Jf(b)) : a = false, tg(b, c), Mg(b, d, e), Og(b, d, e, c), qi(null, b, d, true, a, c);
      case 19:
        return Ai(a, b, c);
      case 23:
        return mi(a, b, c);
      case 24:
        return mi(a, b, c);
    }
    throw Error(y(156, b.tag));
  };
  function ik(a, b, c, d) {
    this.tag = a;
    this.key = c;
    this.sibling = this.child = this.return = this.stateNode = this.type = this.elementType = null;
    this.index = 0;
    this.ref = null;
    this.pendingProps = b;
    this.dependencies = this.memoizedState = this.updateQueue = this.memoizedProps = null;
    this.mode = d;
    this.flags = 0;
    this.lastEffect = this.firstEffect = this.nextEffect = null;
    this.childLanes = this.lanes = 0;
    this.alternate = null;
  }
  function nh(a, b, c, d) {
    return new ik(a, b, c, d);
  }
  function ji(a) {
    a = a.prototype;
    return !(!a || !a.isReactComponent);
  }
  function hk(a) {
    if (typeof a === "function")
      return ji(a) ? 1 : 0;
    if (a !== void 0 && a !== null) {
      a = a.$$typeof;
      if (a === Aa)
        return 11;
      if (a === Da)
        return 14;
    }
    return 2;
  }
  function Tg(a, b) {
    var c = a.alternate;
    c === null ? (c = nh(a.tag, b, a.key, a.mode), c.elementType = a.elementType, c.type = a.type, c.stateNode = a.stateNode, c.alternate = a, a.alternate = c) : (c.pendingProps = b, c.type = a.type, c.flags = 0, c.nextEffect = null, c.firstEffect = null, c.lastEffect = null);
    c.childLanes = a.childLanes;
    c.lanes = a.lanes;
    c.child = a.child;
    c.memoizedProps = a.memoizedProps;
    c.memoizedState = a.memoizedState;
    c.updateQueue = a.updateQueue;
    b = a.dependencies;
    c.dependencies = b === null ? null : {lanes: b.lanes, firstContext: b.firstContext};
    c.sibling = a.sibling;
    c.index = a.index;
    c.ref = a.ref;
    return c;
  }
  function Vg(a, b, c, d, e, f) {
    var g = 2;
    d = a;
    if (typeof a === "function")
      ji(a) && (g = 1);
    else if (typeof a === "string")
      g = 5;
    else
      a:
        switch (a) {
          case ua:
            return Xg(c.children, e, f, b);
          case Ha:
            g = 8;
            e |= 16;
            break;
          case wa:
            g = 8;
            e |= 1;
            break;
          case xa:
            return a = nh(12, c, b, e | 8), a.elementType = xa, a.type = xa, a.lanes = f, a;
          case Ba:
            return a = nh(13, c, b, e), a.type = Ba, a.elementType = Ba, a.lanes = f, a;
          case Ca:
            return a = nh(19, c, b, e), a.elementType = Ca, a.lanes = f, a;
          case Ia:
            return vi(c, e, f, b);
          case Ja:
            return a = nh(24, c, b, e), a.elementType = Ja, a.lanes = f, a;
          default:
            if (typeof a === "object" && a !== null)
              switch (a.$$typeof) {
                case ya:
                  g = 10;
                  break a;
                case za:
                  g = 9;
                  break a;
                case Aa:
                  g = 11;
                  break a;
                case Da:
                  g = 14;
                  break a;
                case Ea:
                  g = 16;
                  d = null;
                  break a;
                case Fa:
                  g = 22;
                  break a;
              }
            throw Error(y(130, a == null ? a : typeof a, ""));
        }
    b = nh(g, c, b, e);
    b.elementType = a;
    b.type = d;
    b.lanes = f;
    return b;
  }
  function Xg(a, b, c, d) {
    a = nh(7, a, d, b);
    a.lanes = c;
    return a;
  }
  function vi(a, b, c, d) {
    a = nh(23, a, d, b);
    a.elementType = Ia;
    a.lanes = c;
    return a;
  }
  function Ug(a, b, c) {
    a = nh(6, a, null, b);
    a.lanes = c;
    return a;
  }
  function Wg(a, b, c) {
    b = nh(4, a.children !== null ? a.children : [], a.key, b);
    b.lanes = c;
    b.stateNode = {containerInfo: a.containerInfo, pendingChildren: null, implementation: a.implementation};
    return b;
  }
  function jk(a, b, c) {
    this.tag = b;
    this.containerInfo = a;
    this.finishedWork = this.pingCache = this.current = this.pendingChildren = null;
    this.timeoutHandle = -1;
    this.pendingContext = this.context = null;
    this.hydrate = c;
    this.callbackNode = null;
    this.callbackPriority = 0;
    this.eventTimes = Zc(0);
    this.expirationTimes = Zc(-1);
    this.entangledLanes = this.finishedLanes = this.mutableReadLanes = this.expiredLanes = this.pingedLanes = this.suspendedLanes = this.pendingLanes = 0;
    this.entanglements = Zc(0);
    this.mutableSourceEagerHydrationData = null;
  }
  function kk(a, b, c) {
    var d = 3 < arguments.length && arguments[3] !== void 0 ? arguments[3] : null;
    return {$$typeof: ta, key: d == null ? null : "" + d, children: a, containerInfo: b, implementation: c};
  }
  function lk(a, b, c, d) {
    var e = b.current, f = Hg(), g = Ig(e);
    a:
      if (c) {
        c = c._reactInternals;
        b: {
          if (Zb(c) !== c || c.tag !== 1)
            throw Error(y(170));
          var h = c;
          do {
            switch (h.tag) {
              case 3:
                h = h.stateNode.context;
                break b;
              case 1:
                if (Ff(h.type)) {
                  h = h.stateNode.__reactInternalMemoizedMergedChildContext;
                  break b;
                }
            }
            h = h.return;
          } while (h !== null);
          throw Error(y(171));
        }
        if (c.tag === 1) {
          var k = c.type;
          if (Ff(k)) {
            c = If(c, k, h);
            break a;
          }
        }
        c = h;
      } else
        c = Cf;
    b.context === null ? b.context = c : b.pendingContext = c;
    b = zg(f, g);
    b.payload = {element: a};
    d = d === void 0 ? null : d;
    d !== null && (b.callback = d);
    Ag(e, b);
    Jg(e, g, f);
    return g;
  }
  function mk(a) {
    a = a.current;
    if (!a.child)
      return null;
    switch (a.child.tag) {
      case 5:
        return a.child.stateNode;
      default:
        return a.child.stateNode;
    }
  }
  function nk(a, b) {
    a = a.memoizedState;
    if (a !== null && a.dehydrated !== null) {
      var c = a.retryLane;
      a.retryLane = c !== 0 && c < b ? c : b;
    }
  }
  function ok(a, b) {
    nk(a, b);
    (a = a.alternate) && nk(a, b);
  }
  function pk() {
    return null;
  }
  function qk(a, b, c) {
    var d = c != null && c.hydrationOptions != null && c.hydrationOptions.mutableSources || null;
    c = new jk(a, b, c != null && c.hydrate === true);
    b = nh(3, null, null, b === 2 ? 7 : b === 1 ? 3 : 0);
    c.current = b;
    b.stateNode = c;
    xg(b);
    a[ff] = c.current;
    cf(a.nodeType === 8 ? a.parentNode : a);
    if (d)
      for (a = 0; a < d.length; a++) {
        b = d[a];
        var e = b._getVersion;
        e = e(b._source);
        c.mutableSourceEagerHydrationData == null ? c.mutableSourceEagerHydrationData = [b, e] : c.mutableSourceEagerHydrationData.push(b, e);
      }
    this._internalRoot = c;
  }
  qk.prototype.render = function(a) {
    lk(a, this._internalRoot, null, null);
  };
  qk.prototype.unmount = function() {
    var a = this._internalRoot, b = a.containerInfo;
    lk(null, a, null, function() {
      b[ff] = null;
    });
  };
  function rk(a) {
    return !(!a || a.nodeType !== 1 && a.nodeType !== 9 && a.nodeType !== 11 && (a.nodeType !== 8 || a.nodeValue !== " react-mount-point-unstable "));
  }
  function sk(a, b) {
    b || (b = a ? a.nodeType === 9 ? a.documentElement : a.firstChild : null, b = !(!b || b.nodeType !== 1 || !b.hasAttribute("data-reactroot")));
    if (!b)
      for (var c; c = a.lastChild; )
        a.removeChild(c);
    return new qk(a, 0, b ? {hydrate: true} : void 0);
  }
  function tk(a, b, c, d, e) {
    var f = c._reactRootContainer;
    if (f) {
      var g = f._internalRoot;
      if (typeof e === "function") {
        var h = e;
        e = function() {
          var a2 = mk(g);
          h.call(a2);
        };
      }
      lk(b, g, a, e);
    } else {
      f = c._reactRootContainer = sk(c, d);
      g = f._internalRoot;
      if (typeof e === "function") {
        var k = e;
        e = function() {
          var a2 = mk(g);
          k.call(a2);
        };
      }
      Xj(function() {
        lk(b, g, a, e);
      });
    }
    return mk(g);
  }
  ec = function(a) {
    if (a.tag === 13) {
      var b = Hg();
      Jg(a, 4, b);
      ok(a, 4);
    }
  };
  fc = function(a) {
    if (a.tag === 13) {
      var b = Hg();
      Jg(a, 67108864, b);
      ok(a, 67108864);
    }
  };
  gc = function(a) {
    if (a.tag === 13) {
      var b = Hg(), c = Ig(a);
      Jg(a, c, b);
      ok(a, c);
    }
  };
  hc = function(a, b) {
    return b();
  };
  yb = function(a, b, c) {
    switch (b) {
      case "input":
        ab(a, c);
        b = c.name;
        if (c.type === "radio" && b != null) {
          for (c = a; c.parentNode; )
            c = c.parentNode;
          c = c.querySelectorAll("input[name=" + JSON.stringify("" + b) + '][type="radio"]');
          for (b = 0; b < c.length; b++) {
            var d = c[b];
            if (d !== a && d.form === a.form) {
              var e = Db(d);
              if (!e)
                throw Error(y(90));
              Wa(d);
              ab(d, e);
            }
          }
        }
        break;
      case "textarea":
        ib(a, c);
        break;
      case "select":
        b = c.value, b != null && fb(a, !!c.multiple, b, false);
    }
  };
  Gb = Wj;
  Hb = function(a, b, c, d, e) {
    var f = X;
    X |= 4;
    try {
      return gg(98, a.bind(null, b, c, d, e));
    } finally {
      X = f, X === 0 && (wj(), ig());
    }
  };
  Ib = function() {
    (X & 49) === 0 && (Vj(), Oj());
  };
  Jb = function(a, b) {
    var c = X;
    X |= 2;
    try {
      return a(b);
    } finally {
      X = c, X === 0 && (wj(), ig());
    }
  };
  function uk(a, b) {
    var c = 2 < arguments.length && arguments[2] !== void 0 ? arguments[2] : null;
    if (!rk(b))
      throw Error(y(200));
    return kk(a, b, null, c);
  }
  var vk = {Events: [Cb, ue, Db, Eb, Fb, Oj, {current: false}]};
  var wk = {findFiberByHostInstance: wc, bundleType: 0, version: "17.0.1", rendererPackageName: "react-dom"};
  var xk = {bundleType: wk.bundleType, version: wk.version, rendererPackageName: wk.rendererPackageName, rendererConfig: wk.rendererConfig, overrideHookState: null, overrideHookStateDeletePath: null, overrideHookStateRenamePath: null, overrideProps: null, overridePropsDeletePath: null, overridePropsRenamePath: null, setSuspenseHandler: null, scheduleUpdate: null, currentDispatcherRef: ra.ReactCurrentDispatcher, findHostInstanceByFiber: function(a) {
    a = cc(a);
    return a === null ? null : a.stateNode;
  }, findFiberByHostInstance: wk.findFiberByHostInstance || pk, findHostInstancesForRefresh: null, scheduleRefresh: null, scheduleRoot: null, setRefreshHandler: null, getCurrentFiber: null};
  if (typeof __REACT_DEVTOOLS_GLOBAL_HOOK__ !== "undefined") {
    yk = __REACT_DEVTOOLS_GLOBAL_HOOK__;
    if (!yk.isDisabled && yk.supportsFiber)
      try {
        Lf = yk.inject(xk), Mf = yk;
      } catch (a) {
      }
  }
  var yk;
  exports2.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED = vk;
  exports2.createPortal = uk;
  exports2.findDOMNode = function(a) {
    if (a == null)
      return null;
    if (a.nodeType === 1)
      return a;
    var b = a._reactInternals;
    if (b === void 0) {
      if (typeof a.render === "function")
        throw Error(y(188));
      throw Error(y(268, Object.keys(a)));
    }
    a = cc(b);
    a = a === null ? null : a.stateNode;
    return a;
  };
  exports2.flushSync = function(a, b) {
    var c = X;
    if ((c & 48) !== 0)
      return a(b);
    X |= 1;
    try {
      if (a)
        return gg(99, a.bind(null, b));
    } finally {
      X = c, ig();
    }
  };
  exports2.hydrate = function(a, b, c) {
    if (!rk(b))
      throw Error(y(200));
    return tk(null, a, b, true, c);
  };
  exports2.render = function(a, b, c) {
    if (!rk(b))
      throw Error(y(200));
    return tk(null, a, b, false, c);
  };
  exports2.unmountComponentAtNode = function(a) {
    if (!rk(a))
      throw Error(y(40));
    return a._reactRootContainer ? (Xj(function() {
      tk(null, null, a, false, function() {
        a._reactRootContainer = null;
        a[ff] = null;
      });
    }), true) : false;
  };
  exports2.unstable_batchedUpdates = Wj;
  exports2.unstable_createPortal = function(a, b) {
    return uk(a, b, 2 < arguments.length && arguments[2] !== void 0 ? arguments[2] : null);
  };
  exports2.unstable_renderSubtreeIntoContainer = function(a, b, c, d) {
    if (!rk(c))
      throw Error(y(200));
    if (a == null || a._reactInternals === void 0)
      throw Error(y(38));
    return tk(a, b, c, false, d);
  };
  exports2.version = "17.0.1";
});

// node_modules/react-dom/index.js
var require_react_dom = __commonJS((exports2, module2) => {
  "use strict";
  function checkDCE() {
    if (typeof __REACT_DEVTOOLS_GLOBAL_HOOK__ === "undefined" || typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE !== "function") {
      return;
    }
    if (false) {
      throw new Error("^_^");
    }
    try {
      __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(checkDCE);
    } catch (err) {
      console.error(err);
    }
  }
  if (true) {
    checkDCE();
    module2.exports = require_react_dom_production_min();
  } else {
    module2.exports = null;
  }
});

// node_modules/clsx/dist/clsx.js
var require_clsx = __commonJS((exports2, module2) => {
  function toVal(mix) {
    var k, y, str = "";
    if (typeof mix === "string" || typeof mix === "number") {
      str += mix;
    } else if (typeof mix === "object") {
      if (Array.isArray(mix)) {
        for (k = 0; k < mix.length; k++) {
          if (mix[k]) {
            if (y = toVal(mix[k])) {
              str && (str += " ");
              str += y;
            }
          }
        }
      } else {
        for (k in mix) {
          if (mix[k]) {
            str && (str += " ");
            str += k;
          }
        }
      }
    }
    return str;
  }
  module2.exports = function() {
    var i = 0, tmp, x, str = "";
    while (i < arguments.length) {
      if (tmp = arguments[i++]) {
        if (x = toVal(tmp)) {
          str && (str += " ");
          str += x;
        }
      }
    }
    return str;
  };
});

// node_modules/react-toastify/dist/react-toastify.cjs.production.min.js
var require_react_toastify_cjs_production_min = __commonJS((exports2) => {
  "use strict";
  function e(e2) {
    return e2 && typeof e2 == "object" && "default" in e2 ? e2.default : e2;
  }
  Object.defineProperty(exports2, "__esModule", {value: true});
  var t = require_react();
  var n = e(t);
  var o = e(require_clsx());
  var r = require_react_dom();
  function s() {
    return (s = Object.assign || function(e2) {
      for (var t2 = 1; t2 < arguments.length; t2++) {
        var n2 = arguments[t2];
        for (var o2 in n2)
          Object.prototype.hasOwnProperty.call(n2, o2) && (e2[o2] = n2[o2]);
      }
      return e2;
    }).apply(this, arguments);
  }
  function a(e2) {
    return typeof e2 == "number" && !isNaN(e2);
  }
  function i(e2) {
    return typeof e2 == "boolean";
  }
  function c(e2) {
    return typeof e2 == "string";
  }
  function u(e2) {
    return typeof e2 == "function";
  }
  function l(e2) {
    return c(e2) || u(e2) ? e2 : null;
  }
  function d(e2) {
    return e2 === 0 || e2;
  }
  var f = !(typeof window == "undefined" || !window.document || !window.document.createElement);
  function p(e2) {
    return t.isValidElement(e2) || c(e2) || u(e2) || a(e2);
  }
  var m = {TOP_LEFT: "top-left", TOP_RIGHT: "top-right", TOP_CENTER: "top-center", BOTTOM_LEFT: "bottom-left", BOTTOM_RIGHT: "bottom-right", BOTTOM_CENTER: "bottom-center"};
  var g = {INFO: "info", SUCCESS: "success", WARNING: "warning", ERROR: "error", DEFAULT: "default", DARK: "dark"};
  function v(e2, t2, n2) {
    n2 === void 0 && (n2 = 300);
    var o2 = e2.scrollHeight, r2 = e2.style;
    requestAnimationFrame(function() {
      r2.minHeight = "initial", r2.height = o2 + "px", r2.transition = "all " + n2 + "ms", requestAnimationFrame(function() {
        r2.height = "0", r2.padding = "0", r2.margin = "0", setTimeout(t2, n2);
      });
    });
  }
  function y(e2) {
    var o2 = e2.enter, r2 = e2.exit, s2 = e2.appendPosition, a2 = s2 !== void 0 && s2, i2 = e2.collapse, c2 = i2 === void 0 || i2, u2 = e2.collapseDuration, l2 = u2 === void 0 ? 300 : u2;
    return function(e3) {
      var s3 = e3.children, i3 = e3.position, u3 = e3.preventExitTransition, d2 = e3.done, f2 = e3.nodeRef, p2 = e3.isIn, m2 = a2 ? o2 + "--" + i3 : o2, g2 = a2 ? r2 + "--" + i3 : r2, y2 = t.useRef(), T2 = t.useRef(0);
      function h2() {
        var e4 = f2.current;
        e4.removeEventListener("animationend", h2), T2.current === 0 && (e4.className = y2.current);
      }
      function b2() {
        var e4 = f2.current;
        e4.removeEventListener("animationend", b2), c2 ? v(e4, d2, l2) : d2();
      }
      return t.useLayoutEffect(function() {
        var e4;
        y2.current = (e4 = f2.current).className, e4.className += " " + m2, e4.addEventListener("animationend", h2);
      }, []), t.useEffect(function() {
        p2 || (u3 ? b2() : function() {
          T2.current = 1;
          var e4 = f2.current;
          e4.className += " " + g2, e4.addEventListener("animationend", b2);
        }());
      }, [p2]), n.createElement(n.Fragment, null, s3);
    };
  }
  var T = {list: new Map(), emitQueue: new Map(), on: function(e2, t2) {
    return this.list.has(e2) || this.list.set(e2, []), this.list.get(e2).push(t2), this;
  }, off: function(e2, t2) {
    if (t2) {
      var n2 = this.list.get(e2).filter(function(e3) {
        return e3 !== t2;
      });
      return this.list.set(e2, n2), this;
    }
    return this.list.delete(e2), this;
  }, cancelEmit: function(e2) {
    var t2 = this.emitQueue.get(e2);
    return t2 && (t2.forEach(clearTimeout), this.emitQueue.delete(e2)), this;
  }, emit: function(e2) {
    for (var t2 = this, n2 = arguments.length, o2 = new Array(n2 > 1 ? n2 - 1 : 0), r2 = 1; r2 < n2; r2++)
      o2[r2 - 1] = arguments[r2];
    this.list.has(e2) && this.list.get(e2).forEach(function(n3) {
      var r3 = setTimeout(function() {
        n3.apply(void 0, o2);
      }, 0);
      t2.emitQueue.has(e2) || t2.emitQueue.set(e2, []), t2.emitQueue.get(e2).push(r3);
    });
  }};
  function h(e2, n2) {
    n2 === void 0 && (n2 = false);
    var o2 = t.useRef(e2);
    return t.useEffect(function() {
      n2 && (o2.current = e2);
    }), o2.current;
  }
  function b(e2, t2) {
    switch (t2.type) {
      case 0:
        return [].concat(e2, [t2.toastId]).filter(function(e3) {
          return e3 !== t2.staleId;
        });
      case 1:
        return d(t2.toastId) ? e2.filter(function(e3) {
          return e3 !== t2.toastId;
        }) : [];
    }
  }
  function E(e2) {
    var n2 = t.useReducer(function(e3) {
      return e3 + 1;
    }, 0)[1], o2 = t.useReducer(b, []), r2 = o2[0], s2 = o2[1], f2 = t.useRef(null), m2 = h(0), g2 = h([]), v2 = h({}), y2 = h({toastKey: 1, displayedToast: 0, props: e2, containerId: null, isToastActive: E2, getToast: function(e3) {
      return v2[e3] || null;
    }});
    function E2(e3) {
      return r2.indexOf(e3) !== -1;
    }
    function O2(e3) {
      var t2 = e3.containerId, n3 = y2.props;
      n3.limit && (!t2 || y2.containerId === t2 && n3.enableMultiContainer) && (m2 -= g2.length, g2 = []);
    }
    function C2(e3) {
      s2({type: 1, toastId: e3});
    }
    function I2() {
      var e3 = g2.shift();
      x2(e3.toastContent, e3.toastProps, e3.staleId);
    }
    function _2(e3, o3) {
      var r3 = o3.delay, s3 = o3.staleId, T2 = function(e4, t2) {
        if (e4 == null)
          return {};
        var n3, o4, r4 = {}, s4 = Object.keys(e4);
        for (o4 = 0; o4 < s4.length; o4++)
          t2.indexOf(n3 = s4[o4]) >= 0 || (r4[n3] = e4[n3]);
        return r4;
      }(o3, ["delay", "staleId"]);
      if (p(e3) && (h2 = T2, !(!f2.current || y2.props.enableMultiContainer && h2.containerId !== y2.props.containerId || v2[h2.toastId] && h2.updateId == null))) {
        var h2, b2 = T2.toastId, E3 = y2.props, O3 = function() {
          return C2(b2);
        }, _3 = T2.updateId == null;
        _3 && m2++;
        var R2, N2, P2 = {toastId: b2, updateId: T2.updateId, isIn: false, key: T2.key || y2.toastKey++, type: T2.type, closeToast: O3, closeButton: T2.closeButton, rtl: E3.rtl, position: T2.position || E3.position, transition: T2.transition || E3.transition, className: l(T2.className || E3.toastClassName), bodyClassName: l(T2.bodyClassName || E3.bodyClassName), style: T2.style || E3.toastStyle, bodyStyle: T2.bodyStyle || E3.bodyStyle, onClick: T2.onClick || E3.onClick, pauseOnHover: i(T2.pauseOnHover) ? T2.pauseOnHover : E3.pauseOnHover, pauseOnFocusLoss: i(T2.pauseOnFocusLoss) ? T2.pauseOnFocusLoss : E3.pauseOnFocusLoss, draggable: i(T2.draggable) ? T2.draggable : E3.draggable, draggablePercent: a(T2.draggablePercent) ? T2.draggablePercent : E3.draggablePercent, draggableDirection: T2.draggableDirection || E3.draggableDirection, closeOnClick: i(T2.closeOnClick) ? T2.closeOnClick : E3.closeOnClick, progressClassName: l(T2.progressClassName || E3.progressClassName), progressStyle: T2.progressStyle || E3.progressStyle, autoClose: (R2 = T2.autoClose, N2 = E3.autoClose, R2 === false || a(R2) && R2 > 0 ? R2 : N2), hideProgressBar: i(T2.hideProgressBar) ? T2.hideProgressBar : E3.hideProgressBar, progress: T2.progress, role: c(T2.role) ? T2.role : E3.role, deleteToast: function() {
          !function(e4) {
            delete v2[e4];
            var t2 = g2.length;
            if ((m2 = d(e4) ? m2 - 1 : m2 - y2.displayedToast) < 0 && (m2 = 0), t2 > 0) {
              var o4 = d(e4) ? 1 : y2.props.limit;
              if (t2 === 1 || o4 === 1)
                y2.displayedToast++, I2();
              else {
                var r4 = o4 > t2 ? t2 : o4;
                y2.displayedToast = r4;
                for (var s4 = 0; s4 < r4; s4++)
                  I2();
              }
            } else
              n2();
          }(b2);
        }};
        u(T2.onOpen) && (P2.onOpen = T2.onOpen), u(T2.onClose) && (P2.onClose = T2.onClose), P2.draggableDirection === "y" && P2.draggablePercent === 80 && (P2.draggablePercent *= 1.5);
        var L2 = E3.closeButton;
        T2.closeButton === false || p(T2.closeButton) ? L2 = T2.closeButton : T2.closeButton === true && (L2 = !p(E3.closeButton) || E3.closeButton), P2.closeButton = L2;
        var k2 = e3;
        t.isValidElement(e3) && !c(e3.type) ? k2 = t.cloneElement(e3, {closeToast: O3, toastProps: P2}) : u(e3) && (k2 = e3({closeToast: O3, toastProps: P2})), E3.limit && E3.limit > 0 && m2 > E3.limit && _3 ? g2.push({toastContent: k2, toastProps: P2, staleId: s3}) : a(r3) && r3 > 0 ? setTimeout(function() {
          x2(k2, P2, s3);
        }, r3) : x2(k2, P2, s3);
      }
    }
    function x2(e3, t2, n3) {
      var o3 = t2.toastId;
      n3 && delete v2[n3], v2[o3] = {content: e3, props: t2}, s2({type: 0, toastId: o3, staleId: n3});
    }
    return t.useEffect(function() {
      return y2.containerId = e2.containerId, T.cancelEmit(3).on(0, _2).on(1, function(e3) {
        return f2.current && C2(e3);
      }).on(5, O2).emit(2, y2), function() {
        return T.emit(3, y2);
      };
    }, []), t.useEffect(function() {
      y2.isToastActive = E2, y2.displayedToast = r2.length, T.emit(4, r2.length, e2.containerId);
    }, [r2]), t.useEffect(function() {
      y2.props = e2;
    }), {getToastToRender: function(t2) {
      for (var n3 = {}, o3 = e2.newestOnTop ? Object.keys(v2).reverse() : Object.keys(v2), r3 = 0; r3 < o3.length; r3++) {
        var s3 = v2[o3[r3]], a2 = s3.props.position;
        n3[a2] || (n3[a2] = []), n3[a2].push(s3);
      }
      return Object.keys(n3).map(function(e3) {
        return t2(e3, n3[e3]);
      });
    }, collection: v2, containerRef: f2, isToastActive: E2};
  }
  function O(e2) {
    return e2.targetTouches && e2.targetTouches.length >= 1 ? e2.targetTouches[0].clientX : e2.clientX;
  }
  function C(e2) {
    return e2.targetTouches && e2.targetTouches.length >= 1 ? e2.targetTouches[0].clientY : e2.clientY;
  }
  function I(e2) {
    var n2 = t.useState(true), o2 = n2[0], r2 = n2[1], s2 = t.useState(false), a2 = s2[0], i2 = s2[1], c2 = t.useRef(null), l2 = h({start: 0, x: 0, y: 0, delta: 0, removalDistance: 0, canCloseOnClick: true, canDrag: false, boundingRect: null}), d2 = h(e2, true), f2 = e2.autoClose, p2 = e2.pauseOnHover, m2 = e2.closeToast, g2 = e2.onClick, v2 = e2.closeOnClick;
    function y2(t2) {
      if (e2.draggable) {
        var n3 = c2.current;
        l2.canCloseOnClick = true, l2.canDrag = true, l2.boundingRect = n3.getBoundingClientRect(), n3.style.transition = "", l2.x = O(t2.nativeEvent), l2.y = C(t2.nativeEvent), e2.draggableDirection === "x" ? (l2.start = l2.x, l2.removalDistance = n3.offsetWidth * (e2.draggablePercent / 100)) : (l2.start = l2.y, l2.removalDistance = n3.offsetHeight * (e2.draggablePercent / 100));
      }
    }
    function T2() {
      if (l2.boundingRect) {
        var t2 = l2.boundingRect;
        e2.pauseOnHover && l2.x >= t2.left && l2.x <= t2.right && l2.y >= t2.top && l2.y <= t2.bottom ? E2() : b2();
      }
    }
    function b2() {
      r2(true);
    }
    function E2() {
      r2(false);
    }
    function I2(t2) {
      if (l2.canDrag) {
        t2.preventDefault();
        var n3 = c2.current;
        o2 && E2(), l2.x = O(t2), l2.y = C(t2), l2.delta = e2.draggableDirection === "x" ? l2.x - l2.start : l2.y - l2.start, l2.start !== l2.x && (l2.canCloseOnClick = false), n3.style.transform = "translate" + e2.draggableDirection + "(" + l2.delta + "px)", n3.style.opacity = "" + (1 - Math.abs(l2.delta / l2.removalDistance));
      }
    }
    function _2() {
      var t2 = c2.current;
      if (l2.canDrag) {
        if (l2.canDrag = false, Math.abs(l2.delta) > l2.removalDistance)
          return i2(true), void e2.closeToast();
        t2.style.transition = "transform 0.2s, opacity 0.2s", t2.style.transform = "translate" + e2.draggableDirection + "(0)", t2.style.opacity = "1";
      }
    }
    t.useEffect(function() {
      return u(e2.onOpen) && e2.onOpen(t.isValidElement(e2.children) && e2.children.props), function() {
        u(d2.onClose) && d2.onClose(t.isValidElement(d2.children) && d2.children.props);
      };
    }, []), t.useEffect(function() {
      return e2.draggable && (document.addEventListener("mousemove", I2), document.addEventListener("mouseup", _2), document.addEventListener("touchmove", I2), document.addEventListener("touchend", _2)), function() {
        e2.draggable && (document.removeEventListener("mousemove", I2), document.removeEventListener("mouseup", _2), document.removeEventListener("touchmove", I2), document.removeEventListener("touchend", _2));
      };
    }, [e2.draggable]), t.useEffect(function() {
      return e2.pauseOnFocusLoss && (document.hasFocus() || E2(), window.addEventListener("focus", b2), window.addEventListener("blur", E2)), function() {
        e2.pauseOnFocusLoss && (window.removeEventListener("focus", b2), window.removeEventListener("blur", E2));
      };
    }, [e2.pauseOnFocusLoss]);
    var x2 = {onMouseDown: y2, onTouchStart: y2, onMouseUp: T2, onTouchEnd: T2};
    return f2 && p2 && (x2.onMouseEnter = E2, x2.onMouseLeave = b2), v2 && (x2.onClick = function(e3) {
      g2 && g2(e3), l2.canCloseOnClick && m2();
    }), {playToast: b2, pauseToast: E2, isRunning: o2, preventExitTransition: a2, toastRef: c2, eventHandlers: x2};
  }
  function _(e2) {
    var n2 = e2.closeToast, o2 = e2.ariaLabel;
    return t.createElement("button", {className: "Toastify__close-button Toastify__close-button--" + e2.type, type: "button", onClick: function(e3) {
      e3.stopPropagation(), n2(e3);
    }, "aria-label": o2 === void 0 ? "close" : o2}, t.createElement("svg", {"aria-hidden": "true", viewBox: "0 0 14 16"}, t.createElement("path", {fillRule: "evenodd", d: "M7.71 8.23l3.75 3.75-1.48 1.48-3.75-3.75-3.75 3.75L1 11.98l3.75-3.75L1 4.48 2.48 3l3.75 3.75L9.98 3l1.48 1.48-3.75 3.75z"})));
  }
  function x(e2) {
    var n2, r2, a2 = e2.closeToast, i2 = e2.type, c2 = e2.className, l2 = e2.controlledProgress, d2 = e2.progress, f2 = e2.rtl, p2 = e2.isIn, m2 = s({}, e2.style, {animationDuration: e2.delay + "ms", animationPlayState: e2.isRunning ? "running" : "paused", opacity: e2.hide ? 0 : 1});
    l2 && (m2.transform = "scaleX(" + d2 + ")");
    var g2 = o("Toastify__progress-bar", l2 ? "Toastify__progress-bar--controlled" : "Toastify__progress-bar--animated", "Toastify__progress-bar--" + i2, ((n2 = {})["Toastify__progress-bar--rtl"] = f2, n2)), v2 = u(c2) ? c2({rtl: f2, type: i2, defaultClassName: g2}) : o(g2, c2), y2 = ((r2 = {})[l2 && d2 >= 1 ? "onTransitionEnd" : "onAnimationEnd"] = l2 && d2 < 1 ? null : function() {
      p2 && a2();
    }, r2);
    return t.createElement("div", Object.assign({role: "progressbar", className: v2, style: m2}, y2));
  }
  x.defaultProps = {type: g.DEFAULT, hide: false};
  var R = function(e2) {
    var n2, r2 = I(e2), s2 = r2.isRunning, a2 = r2.preventExitTransition, i2 = r2.toastRef, c2 = r2.eventHandlers, l2 = e2.closeButton, d2 = e2.children, f2 = e2.autoClose, p2 = e2.onClick, m2 = e2.type, g2 = e2.hideProgressBar, v2 = e2.closeToast, y2 = e2.transition, T2 = e2.position, h2 = e2.className, b2 = e2.style, E2 = e2.bodyClassName, O2 = e2.bodyStyle, C2 = e2.progressClassName, _2 = e2.progressStyle, R2 = e2.updateId, N2 = e2.role, P2 = e2.progress, L2 = e2.rtl, k2 = e2.toastId, w2 = e2.deleteToast, D2 = e2.isIn, B2 = o("Toastify__toast", "Toastify__toast--" + m2, ((n2 = {})["Toastify__toast--rtl"] = L2, n2)), S2 = u(h2) ? h2({rtl: L2, position: T2, type: m2, defaultClassName: B2}) : o(B2, h2), F2 = !!P2;
    return t.createElement(y2, {isIn: D2, done: w2, position: T2, preventExitTransition: a2, nodeRef: i2}, t.createElement("div", Object.assign({id: k2, onClick: p2, className: S2}, c2, {style: b2, ref: i2}), t.createElement("div", Object.assign({}, D2 && {role: N2}, {className: u(E2) ? E2({type: m2}) : o("Toastify__toast-body", E2), style: O2}), d2), function(e3) {
      if (e3) {
        var n3 = {closeToast: v2, type: m2};
        return u(e3) ? e3(n3) : t.isValidElement(e3) ? t.cloneElement(e3, n3) : void 0;
      }
    }(l2), (f2 || F2) && t.createElement(x, Object.assign({}, R2 && !F2 ? {key: "pb-" + R2} : {}, {rtl: L2, delay: f2, isRunning: s2, isIn: D2, closeToast: v2, hide: g2, type: m2, style: _2, className: C2, controlledProgress: F2, progress: P2}))));
  };
  var N = y({enter: "Toastify--animate Toastify__bounce-enter", exit: "Toastify--animate Toastify__bounce-exit", appendPosition: true});
  var P = y({enter: "Toastify--animate Toastify__slide-enter", exit: "Toastify--animate Toastify__slide-exit", appendPosition: true});
  var L = y({enter: "Toastify--animate Toastify__zoom-enter", exit: "Toastify--animate Toastify__zoom-exit"});
  var k = y({enter: "Toastify--animate Toastify__flip-enter", exit: "Toastify--animate Toastify__flip-exit"});
  var w = function(e2) {
    var n2 = E(e2), r2 = n2.isToastActive, a2 = e2.className, i2 = e2.style, c2 = e2.rtl;
    function d2(e3) {
      var t2, n3 = o("Toastify__toast-container", "Toastify__toast-container--" + e3, ((t2 = {})["Toastify__toast-container--rtl"] = c2, t2));
      return u(a2) ? a2({position: e3, rtl: c2, defaultClassName: n3}) : o(n3, l(a2));
    }
    return t.createElement("div", {ref: n2.containerRef, className: "Toastify", id: e2.containerId}, (0, n2.getToastToRender)(function(e3, n3) {
      var o2 = n3.length === 0 ? s({}, i2, {pointerEvents: "none"}) : s({}, i2);
      return t.createElement("div", {className: d2(e3), style: o2, key: "container-" + e3}, n3.map(function(e4) {
        var n4 = e4.content, o3 = e4.props;
        return t.createElement(R, Object.assign({}, o3, {isIn: r2(o3.toastId), key: "toast-" + o3.key, closeButton: o3.closeButton === true ? _ : o3.closeButton}), n4);
      }));
    }));
  };
  w.defaultProps = {position: m.TOP_RIGHT, transition: N, rtl: false, autoClose: 5e3, hideProgressBar: false, closeButton: _, pauseOnHover: true, pauseOnFocusLoss: true, closeOnClick: true, newestOnTop: false, draggable: true, draggablePercent: 80, draggableDirection: "x", role: "alert"};
  var D;
  var B;
  var S;
  var F = new Map();
  var A = [];
  var M = false;
  function H() {
    return Math.random().toString(36).substr(2, 9);
  }
  function j(e2) {
    return e2 && (c(e2.toastId) || a(e2.toastId)) ? e2.toastId : H();
  }
  function Q(e2, n2) {
    return F.size > 0 ? T.emit(0, e2, n2) : (A.push({content: e2, options: n2}), M && f && (M = false, B = document.createElement("div"), document.body.appendChild(B), r.render(t.createElement(w, Object.assign({}, S)), B))), n2.toastId;
  }
  function U(e2, t2) {
    return s({}, t2, {type: t2 && t2.type || e2, toastId: j(t2)});
  }
  var q = function(e2) {
    return function(t2, n2) {
      return Q(t2, U(e2, n2));
    };
  };
  var z = function(e2, t2) {
    return Q(e2, U(g.DEFAULT, t2));
  };
  z.success = q(g.SUCCESS), z.info = q(g.INFO), z.error = q(g.ERROR), z.warning = q(g.WARNING), z.dark = q(g.DARK), z.warn = z.warning, z.dismiss = function(e2) {
    return T.emit(1, e2);
  }, z.clearWaitingQueue = function(e2) {
    return e2 === void 0 && (e2 = {}), T.emit(5, e2);
  }, z.isActive = function(e2) {
    var t2 = false;
    return F.forEach(function(n2) {
      n2.isToastActive && n2.isToastActive(e2) && (t2 = true);
    }), t2;
  }, z.update = function(e2, t2) {
    t2 === void 0 && (t2 = {}), setTimeout(function() {
      var n2 = function(e3, t3) {
        var n3 = F.get(t3.containerId || D);
        return n3 ? n3.getToast(e3) : null;
      }(e2, t2);
      if (n2) {
        var o2 = n2.content, r2 = s({}, n2.props, t2, {toastId: t2.toastId || e2, updateId: H()});
        r2.toastId !== e2 && (r2.staleId = e2);
        var a2 = r2.render || o2;
        delete r2.render, Q(a2, r2);
      }
    }, 0);
  }, z.done = function(e2) {
    z.update(e2, {progress: 1});
  }, z.onChange = function(e2) {
    return u(e2) && T.on(4, e2), function() {
      u(e2) && T.off(4, e2);
    };
  }, z.configure = function(e2) {
    e2 === void 0 && (e2 = {}), M = true, S = e2;
  }, z.POSITION = m, z.TYPE = g, T.on(2, function(e2) {
    F.set(D = e2.containerId || e2, e2), A.forEach(function(e3) {
      T.emit(0, e3.content, e3.options);
    }), A = [];
  }).on(3, function(e2) {
    F.delete(e2.containerId || e2), F.size === 0 && T.off(0).off(1).off(5), f && B && document.body.removeChild(B);
  }), exports2.Bounce = N, exports2.Flip = k, exports2.Slide = P, exports2.ToastContainer = w, exports2.Zoom = L, exports2.collapseToast = v, exports2.cssTransition = y, exports2.toast = z, exports2.useToast = I, exports2.useToastContainer = E;
});

// node_modules/react-toastify/dist/index.js
var require_dist = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_react_toastify_cjs_production_min();
  } else {
    module2.exports = null;
  }
});

// public/services/toastify.tsx
var require_toastify = __commonJS((exports2) => {
  __markAsModule(exports2);
  __export(exports2, {
    error: () => error2,
    success: () => success2
  });
  var import_react51 = __toModule(require_react());
  var import_react_dom2 = __toModule(require_react_dom());
  var import_react_toastify = __toModule(require_dist());
  var hasContainer = false;
  var setup = () => {
    if (!hasContainer) {
      hasContainer = true;
      import_react_dom2.default.render(/* @__PURE__ */ import_react51.default.createElement(import_react_toastify.ToastContainer, {
        position: import_react_toastify.toast.POSITION.TOP_RIGHT,
        toastClassName: "c-toast"
      }), document.getElementById("root-toastify"));
    }
  };
  var success2 = (content, options) => {
    setup();
    import_react_toastify.toast.success(content, options);
  };
  var error2 = (content, options) => {
    setup();
    import_react_toastify.toast.error(content, options);
  };
});

// node_modules/react-icons/lib/cjs/iconsManifest.js
var require_iconsManifest = __commonJS((exports2, module2) => {
  module2.exports.IconsManifest = [
    {
      id: "fa",
      name: "Font Awesome",
      projectUrl: "https://fontawesome.com/",
      license: "CC BY 4.0 License",
      licenseUrl: "https://creativecommons.org/licenses/by/4.0/"
    },
    {
      id: "io",
      name: "Ionicons 4",
      projectUrl: "https://ionicons.com/",
      license: "MIT",
      licenseUrl: "https://github.com/ionic-team/ionicons/blob/master/LICENSE"
    },
    {
      id: "io5",
      name: "Ionicons 5",
      projectUrl: "https://ionicons.com/",
      license: "MIT",
      licenseUrl: "https://github.com/ionic-team/ionicons/blob/master/LICENSE"
    },
    {
      id: "md",
      name: "Material Design icons",
      projectUrl: "http://google.github.io/material-design-icons/",
      license: "Apache License Version 2.0",
      licenseUrl: "https://github.com/google/material-design-icons/blob/master/LICENSE"
    },
    {
      id: "ti",
      name: "Typicons",
      projectUrl: "http://s-ings.com/typicons/",
      license: "CC BY-SA 3.0",
      licenseUrl: "https://creativecommons.org/licenses/by-sa/3.0/"
    },
    {
      id: "go",
      name: "Github Octicons icons",
      projectUrl: "https://octicons.github.com/",
      license: "MIT",
      licenseUrl: "https://github.com/primer/octicons/blob/master/LICENSE"
    },
    {
      id: "fi",
      name: "Feather",
      projectUrl: "https://feathericons.com/",
      license: "MIT",
      licenseUrl: "https://github.com/feathericons/feather/blob/master/LICENSE"
    },
    {
      id: "gi",
      name: "Game Icons",
      projectUrl: "https://game-icons.net/",
      license: "CC BY 3.0",
      licenseUrl: "https://creativecommons.org/licenses/by/3.0/"
    },
    {
      id: "wi",
      name: "Weather Icons",
      projectUrl: "https://erikflowers.github.io/weather-icons/",
      license: "SIL OFL 1.1",
      licenseUrl: "http://scripts.sil.org/OFL"
    },
    {
      id: "di",
      name: "Devicons",
      projectUrl: "https://vorillaz.github.io/devicons/",
      license: "MIT",
      licenseUrl: "https://opensource.org/licenses/MIT"
    },
    {
      id: "ai",
      name: "Ant Design Icons",
      projectUrl: "https://github.com/ant-design/ant-design-icons",
      license: "MIT",
      licenseUrl: "https://opensource.org/licenses/MIT"
    },
    {
      id: "bs",
      name: "Bootstrap Icons",
      projectUrl: "https://github.com/twbs/icons",
      license: "MIT",
      licenseUrl: "https://opensource.org/licenses/MIT"
    },
    {
      id: "ri",
      name: "Remix Icon",
      projectUrl: "https://github.com/Remix-Design/RemixIcon",
      license: "Apache License Version 2.0",
      licenseUrl: "http://www.apache.org/licenses/"
    },
    {
      id: "fc",
      name: "Flat Color Icons",
      projectUrl: "https://github.com/icons8/flat-color-icons",
      license: "MIT",
      licenseUrl: "https://opensource.org/licenses/MIT"
    },
    {
      id: "gr",
      name: "Grommet-Icons",
      projectUrl: "https://github.com/grommet/grommet-icons",
      license: "Apache License Version 2.0",
      licenseUrl: "http://www.apache.org/licenses/"
    },
    {
      id: "hi",
      name: "Heroicons",
      projectUrl: "https://github.com/tailwindlabs/heroicons",
      license: "MIT",
      licenseUrl: "https://opensource.org/licenses/MIT"
    },
    {
      id: "si",
      name: "Simple Icons",
      projectUrl: "https://simpleicons.org/",
      license: "CC0 1.0 Universal",
      licenseUrl: "https://creativecommons.org/publicdomain/zero/1.0/"
    },
    {
      id: "im",
      name: "IcoMoon Free",
      projectUrl: "https://github.com/Keyamoon/IcoMoon-Free",
      license: "CC BY 4.0 License"
    },
    {
      id: "bi",
      name: "BoxIcons",
      projectUrl: "https://github.com/atisawd/boxicons",
      license: "CC BY 4.0 License"
    },
    {
      id: "cg",
      name: "css.gg",
      projectUrl: "https://github.com/astrit/css.gg",
      license: "MIT",
      licenseUrl: "https://opensource.org/licenses/MIT"
    },
    {
      id: "vsc",
      name: "VS Code Icons",
      projectUrl: "https://github.com/microsoft/vscode-codicons",
      license: "CC BY 4.0",
      licenseUrl: "https://creativecommons.org/licenses/by/4.0/"
    }
  ];
});

// node_modules/react-icons/lib/cjs/iconContext.js
var require_iconContext = __commonJS((exports2) => {
  "use strict";
  Object.defineProperty(exports2, "__esModule", {value: true});
  exports2.IconContext = exports2.DefaultContext = void 0;
  var React48 = require_react();
  exports2.DefaultContext = {
    color: void 0,
    size: void 0,
    className: void 0,
    style: void 0,
    attr: void 0
  };
  exports2.IconContext = React48.createContext && React48.createContext(exports2.DefaultContext);
});

// node_modules/react-icons/lib/cjs/iconBase.js
var require_iconBase = __commonJS((exports2) => {
  "use strict";
  var __assign = exports2 && exports2.__assign || function() {
    __assign = Object.assign || function(t) {
      for (var s, i = 1, n = arguments.length; i < n; i++) {
        s = arguments[i];
        for (var p in s)
          if (Object.prototype.hasOwnProperty.call(s, p))
            t[p] = s[p];
      }
      return t;
    };
    return __assign.apply(this, arguments);
  };
  var __rest = exports2 && exports2.__rest || function(s, e) {
    var t = {};
    for (var p in s)
      if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
      for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
        if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
          t[p[i]] = s[p[i]];
      }
    return t;
  };
  Object.defineProperty(exports2, "__esModule", {value: true});
  exports2.IconBase = exports2.GenIcon = void 0;
  var React48 = require_react();
  var iconContext_1 = require_iconContext();
  function Tree2Element(tree) {
    return tree && tree.map(function(node, i) {
      return React48.createElement(node.tag, __assign({key: i}, node.attr), Tree2Element(node.child));
    });
  }
  function GenIcon2(data) {
    return function(props) {
      return React48.createElement(IconBase, __assign({attr: __assign({}, data.attr)}, props), Tree2Element(data.child));
    };
  }
  exports2.GenIcon = GenIcon2;
  function IconBase(props) {
    var elem = function(conf) {
      var attr = props.attr, size = props.size, title = props.title, svgProps = __rest(props, ["attr", "size", "title"]);
      var computedSize = size || conf.size || "1em";
      var className;
      if (conf.className)
        className = conf.className;
      if (props.className)
        className = (className ? className + " " : "") + props.className;
      return React48.createElement("svg", __assign({stroke: "currentColor", fill: "currentColor", strokeWidth: "0"}, conf.attr, attr, svgProps, {className, style: __assign(__assign({color: props.color || conf.color}, conf.style), props.style), height: computedSize, width: computedSize, xmlns: "http://www.w3.org/2000/svg"}), title && React48.createElement("title", null, title), props.children);
    };
    return iconContext_1.IconContext !== void 0 ? React48.createElement(iconContext_1.IconContext.Consumer, null, function(conf) {
      return elem(conf);
    }) : elem(iconContext_1.DefaultContext);
  }
  exports2.IconBase = IconBase;
});

// node_modules/react-icons/lib/cjs/index.js
var require_cjs = __commonJS((exports2) => {
  "use strict";
  var __createBinding = exports2 && exports2.__createBinding || (Object.create ? function(o, m, k, k2) {
    if (k2 === void 0)
      k2 = k;
    Object.defineProperty(o, k2, {enumerable: true, get: function() {
      return m[k];
    }});
  } : function(o, m, k, k2) {
    if (k2 === void 0)
      k2 = k;
    o[k2] = m[k];
  });
  var __exportStar2 = exports2 && exports2.__exportStar || function(m, exports3) {
    for (var p in m)
      if (p !== "default" && !exports3.hasOwnProperty(p))
        __createBinding(exports3, m, p);
  };
  Object.defineProperty(exports2, "__esModule", {value: true});
  __exportStar2(require_iconsManifest(), exports2);
  __exportStar2(require_iconBase(), exports2);
  __exportStar2(require_iconContext(), exports2);
});

// node_modules/@babel/runtime/helpers/extends.js
var require_extends = __commonJS((exports2, module2) => {
  function _extends() {
    module2.exports = _extends = Object.assign || function(target) {
      for (var i = 1; i < arguments.length; i++) {
        var source = arguments[i];
        for (var key in source) {
          if (Object.prototype.hasOwnProperty.call(source, key)) {
            target[key] = source[key];
          }
        }
      }
      return target;
    };
    module2.exports["default"] = module2.exports, module2.exports.__esModule = true;
    return _extends.apply(this, arguments);
  }
  module2.exports = _extends;
  module2.exports["default"] = module2.exports, module2.exports.__esModule = true;
});

// node_modules/@babel/runtime/helpers/objectWithoutPropertiesLoose.js
var require_objectWithoutPropertiesLoose = __commonJS((exports2, module2) => {
  function _objectWithoutPropertiesLoose(source, excluded) {
    if (source == null)
      return {};
    var target = {};
    var sourceKeys = Object.keys(source);
    var key, i;
    for (i = 0; i < sourceKeys.length; i++) {
      key = sourceKeys[i];
      if (excluded.indexOf(key) >= 0)
        continue;
      target[key] = source[key];
    }
    return target;
  }
  module2.exports = _objectWithoutPropertiesLoose;
  module2.exports["default"] = module2.exports, module2.exports.__esModule = true;
});

// node_modules/use-isomorphic-layout-effect/dist/use-isomorphic-layout-effect.cjs.prod.js
var require_use_isomorphic_layout_effect_cjs_prod = __commonJS((exports2) => {
  "use strict";
  Object.defineProperty(exports2, "__esModule", {
    value: true
  });
  var react = require_react();
  var index = typeof document != "undefined" ? react.useLayoutEffect : react.useEffect;
  exports2.default = index;
});

// node_modules/use-isomorphic-layout-effect/dist/use-isomorphic-layout-effect.cjs.js
var require_use_isomorphic_layout_effect_cjs = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_use_isomorphic_layout_effect_cjs_prod();
  } else {
    module2.exports = null;
  }
});

// node_modules/use-latest/dist/use-latest.cjs.prod.js
var require_use_latest_cjs_prod = __commonJS((exports2) => {
  "use strict";
  function _interopDefault(ex) {
    return ex && typeof ex == "object" && "default" in ex ? ex.default : ex;
  }
  Object.defineProperty(exports2, "__esModule", {
    value: true
  });
  var React48 = require_react();
  var useIsomorphicLayoutEffect = _interopDefault(require_use_isomorphic_layout_effect_cjs());
  var useLatest = function(value) {
    var ref = React48.useRef(value);
    return useIsomorphicLayoutEffect(function() {
      ref.current = value;
    }), ref;
  };
  exports2.default = useLatest;
});

// node_modules/use-latest/dist/use-latest.cjs.js
var require_use_latest_cjs = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_use_latest_cjs_prod();
  } else {
    module2.exports = null;
  }
});

// node_modules/use-composed-ref/dist/use-composed-ref.cjs.js
var require_use_composed_ref_cjs = __commonJS((exports2) => {
  "use strict";
  Object.defineProperty(exports2, "__esModule", {value: true});
  var React48 = require_react();
  var updateRef = function updateRef2(ref, value) {
    if (typeof ref === "function") {
      ref(value);
      return;
    }
    ref.current = value;
  };
  var useComposedRef = function useComposedRef2(libRef, userRef) {
    var prevUserRef = React48.useRef();
    return React48.useCallback(function(instance) {
      libRef.current = instance;
      if (prevUserRef.current) {
        updateRef(prevUserRef.current, null);
      }
      prevUserRef.current = userRef;
      if (!userRef) {
        return;
      }
      updateRef(userRef, instance);
    }, [userRef]);
  };
  exports2.default = useComposedRef;
});

// node_modules/react-textarea-autosize/dist/react-textarea-autosize.cjs.prod.js
var require_react_textarea_autosize_cjs_prod = __commonJS((exports2) => {
  "use strict";
  Object.defineProperty(exports2, "__esModule", {value: true});
  var _extends = require_extends();
  var _objectWithoutPropertiesLoose = require_objectWithoutPropertiesLoose();
  var React48 = require_react();
  var useLatest = require_use_latest_cjs();
  var useComposedRef = require_use_composed_ref_cjs();
  function _interopDefault(e) {
    return e && e.__esModule ? e : {default: e};
  }
  var useLatest__default = /* @__PURE__ */ _interopDefault(useLatest);
  var useComposedRef__default = /* @__PURE__ */ _interopDefault(useComposedRef);
  var HIDDEN_TEXTAREA_STYLE = {
    "min-height": "0",
    "max-height": "none",
    height: "0",
    visibility: "hidden",
    overflow: "hidden",
    position: "absolute",
    "z-index": "-1000",
    top: "0",
    right: "0"
  };
  var forceHiddenStyles = function forceHiddenStyles2(node) {
    Object.keys(HIDDEN_TEXTAREA_STYLE).forEach(function(key) {
      node.style.setProperty(key, HIDDEN_TEXTAREA_STYLE[key], "important");
    });
  };
  var hiddenTextarea = null;
  var getHeight = function getHeight2(node, sizingData) {
    var height = node.scrollHeight;
    if (sizingData.sizingStyle.boxSizing === "border-box") {
      return height + sizingData.borderSize;
    }
    return height - sizingData.paddingSize;
  };
  function calculateNodeHeight(sizingData, value, minRows, maxRows) {
    if (minRows === void 0) {
      minRows = 1;
    }
    if (maxRows === void 0) {
      maxRows = Infinity;
    }
    if (!hiddenTextarea) {
      hiddenTextarea = document.createElement("textarea");
      hiddenTextarea.setAttribute("tab-index", "-1");
      hiddenTextarea.setAttribute("aria-hidden", "true");
      forceHiddenStyles(hiddenTextarea);
    }
    if (hiddenTextarea.parentNode === null) {
      document.body.appendChild(hiddenTextarea);
    }
    var paddingSize = sizingData.paddingSize, borderSize = sizingData.borderSize, sizingStyle = sizingData.sizingStyle;
    var boxSizing = sizingStyle.boxSizing;
    Object.keys(sizingStyle).forEach(function(_key) {
      var key = _key;
      hiddenTextarea.style[key] = sizingStyle[key];
    });
    forceHiddenStyles(hiddenTextarea);
    hiddenTextarea.value = value;
    var height = getHeight(hiddenTextarea, sizingData);
    hiddenTextarea.value = "x";
    var rowHeight = hiddenTextarea.scrollHeight - paddingSize;
    var minHeight = rowHeight * minRows;
    if (boxSizing === "border-box") {
      minHeight = minHeight + paddingSize + borderSize;
    }
    height = Math.max(minHeight, height);
    var maxHeight = rowHeight * maxRows;
    if (boxSizing === "border-box") {
      maxHeight = maxHeight + paddingSize + borderSize;
    }
    height = Math.min(maxHeight, height);
    return [height, rowHeight];
  }
  var noop = function noop2() {
  };
  var pick = function pick2(props, obj) {
    return props.reduce(function(acc, prop) {
      acc[prop] = obj[prop];
      return acc;
    }, {});
  };
  var SIZING_STYLE = [
    "borderBottomWidth",
    "borderLeftWidth",
    "borderRightWidth",
    "borderTopWidth",
    "boxSizing",
    "fontFamily",
    "fontSize",
    "fontStyle",
    "fontWeight",
    "letterSpacing",
    "lineHeight",
    "paddingBottom",
    "paddingLeft",
    "paddingRight",
    "paddingTop",
    "tabSize",
    "textIndent",
    "textRendering",
    "textTransform",
    "width"
  ];
  var isIE = typeof document !== "undefined" ? !!document.documentElement.currentStyle : false;
  var getSizingData = function getSizingData2(node) {
    var style = window.getComputedStyle(node);
    if (style === null) {
      return null;
    }
    var sizingStyle = pick(SIZING_STYLE, style);
    var boxSizing = sizingStyle.boxSizing;
    if (boxSizing === "") {
      return null;
    }
    if (isIE && boxSizing === "border-box") {
      sizingStyle.width = parseFloat(sizingStyle.width) + parseFloat(sizingStyle.borderRightWidth) + parseFloat(sizingStyle.borderLeftWidth) + parseFloat(sizingStyle.paddingRight) + parseFloat(sizingStyle.paddingLeft) + "px";
    }
    var paddingSize = parseFloat(sizingStyle.paddingBottom) + parseFloat(sizingStyle.paddingTop);
    var borderSize = parseFloat(sizingStyle.borderBottomWidth) + parseFloat(sizingStyle.borderTopWidth);
    return {
      sizingStyle,
      paddingSize,
      borderSize
    };
  };
  var useWindowResizeListener = function useWindowResizeListener2(listener) {
    var latestListener = useLatest__default["default"](listener);
    React48.useLayoutEffect(function() {
      var handler = function handler2(event) {
        latestListener.current(event);
      };
      window.addEventListener("resize", handler);
      return function() {
        window.removeEventListener("resize", handler);
      };
    }, []);
  };
  var TextareaAutosize = function TextareaAutosize2(_ref, userRef) {
    var cacheMeasurements = _ref.cacheMeasurements, maxRows = _ref.maxRows, minRows = _ref.minRows, _ref$onChange = _ref.onChange, onChange = _ref$onChange === void 0 ? noop : _ref$onChange, _ref$onHeightChange = _ref.onHeightChange, onHeightChange = _ref$onHeightChange === void 0 ? noop : _ref$onHeightChange, props = _objectWithoutPropertiesLoose(_ref, ["cacheMeasurements", "maxRows", "minRows", "onChange", "onHeightChange"]);
    var isControlled = props.value !== void 0;
    var libRef = React48.useRef(null);
    var ref = useComposedRef__default["default"](libRef, userRef);
    var heightRef = React48.useRef(0);
    var measurementsCacheRef = React48.useRef();
    var resizeTextarea = function resizeTextarea2() {
      var node = libRef.current;
      var nodeSizingData = cacheMeasurements && measurementsCacheRef.current ? measurementsCacheRef.current : getSizingData(node);
      if (!nodeSizingData) {
        return;
      }
      measurementsCacheRef.current = nodeSizingData;
      var _calculateNodeHeight = calculateNodeHeight(nodeSizingData, node.value || node.placeholder || "x", minRows, maxRows), height = _calculateNodeHeight[0], rowHeight = _calculateNodeHeight[1];
      if (heightRef.current !== height) {
        heightRef.current = height;
        node.style.setProperty("height", height + "px", "important");
        onHeightChange(height, {
          rowHeight
        });
      }
    };
    var handleChange = function handleChange2(event) {
      if (!isControlled) {
        resizeTextarea();
      }
      onChange(event);
    };
    if (typeof document !== "undefined") {
      React48.useLayoutEffect(resizeTextarea);
      useWindowResizeListener(resizeTextarea);
    }
    return /* @__PURE__ */ React48.createElement("textarea", _extends({}, props, {
      onChange: handleChange,
      ref
    }));
  };
  var index = /* @__PURE__ */ React48.forwardRef(TextareaAutosize);
  exports2.default = index;
});

// node_modules/react-textarea-autosize/dist/react-textarea-autosize.cjs.js
var require_react_textarea_autosize_cjs = __commonJS((exports2, module2) => {
  "use strict";
  if (true) {
    module2.exports = require_react_textarea_autosize_cjs_prod();
  } else {
    module2.exports = null;
  }
});

// node_modules/react-icons/index.js
var require_react_icons = __commonJS((exports2, module2) => {
  module2.exports = require_cjs();
});

// public/ssr.tsx
__markAsModule(exports);
__export(exports, {
  doWork: () => doWork
});
var import_react50 = __toModule(require_react());
var import_server = __toModule(require_server());

// public/pages/Home/Home.page.tsx
var import_react49 = __toModule(require_react());

// public/models/post.ts
var _PostStatus = class {
  constructor(title, value, show, closed, filterable) {
    this.title = title;
    this.value = value;
    this.show = show;
    this.closed = closed;
    this.filterable = filterable;
  }
  static Get(value) {
    for (const status of _PostStatus.All) {
      if (status.value === value) {
        return status;
      }
    }
    throw new Error(`PostStatus not found for value ${value}.`);
  }
};
var PostStatus = _PostStatus;
PostStatus.Open = new _PostStatus("Open", "open", false, false, false);
PostStatus.Planned = new _PostStatus("Planned", "planned", true, false, true);
PostStatus.Started = new _PostStatus("Started", "started", true, false, true);
PostStatus.Completed = new _PostStatus("Completed", "completed", true, true, true);
PostStatus.Declined = new _PostStatus("Declined", "declined", true, true, true);
PostStatus.Duplicate = new _PostStatus("Duplicate", "duplicate", true, true, false);
PostStatus.Deleted = new _PostStatus("Deleted", "deleted", false, true, false);
PostStatus.All = [_PostStatus.Open, _PostStatus.Planned, _PostStatus.Started, _PostStatus.Completed, _PostStatus.Duplicate, _PostStatus.Declined];

// public/models/identity.ts
var TenantStatus;
(function(TenantStatus2) {
  TenantStatus2[TenantStatus2["Active"] = 1] = "Active";
  TenantStatus2[TenantStatus2["Pending"] = 2] = "Pending";
  TenantStatus2[TenantStatus2["Locked"] = 3] = "Locked";
})(TenantStatus || (TenantStatus = {}));
var UserAvatarType;
(function(UserAvatarType2) {
  UserAvatarType2["Letter"] = "letter";
  UserAvatarType2["Gravatar"] = "gravatar";
  UserAvatarType2["Custom"] = "custom";
})(UserAvatarType || (UserAvatarType = {}));
var UserStatus;
(function(UserStatus2) {
  UserStatus2["Active"] = "active";
  UserStatus2["Deleted"] = "deleted";
  UserStatus2["Blocked"] = "blocked";
})(UserStatus || (UserStatus = {}));
var UserRole;
(function(UserRole4) {
  UserRole4["Visitor"] = "visitor";
  UserRole4["Collaborator"] = "collaborator";
  UserRole4["Administrator"] = "administrator";
})(UserRole || (UserRole = {}));
var isCollaborator = (role) => {
  return role === UserRole.Collaborator || role === UserRole.Administrator;
};

// public/components/ErrorBoundary.tsx
var import_react5 = __toModule(require_react());

// public/pages/Error/Error.page.tsx
var import_react4 = __toModule(require_react());

// public/hooks/use-timeout.ts
var import_react = __toModule(require_react());
function useTimeout(callback, delay) {
  const savedCallback = (0, import_react.useRef)();
  (0, import_react.useEffect)(() => {
    savedCallback.current = callback;
  });
  (0, import_react.useEffect)(() => {
    function tick() {
      if (savedCallback.current) {
        savedCallback.current();
      }
    }
    const timer = window.setTimeout(tick, delay);
    return function cleanup() {
      window.clearTimeout(timer);
    };
  }, [delay]);
}

// public/hooks/use-fider.ts
var import_react3 = __toModule(require_react());

// public/services/http.ts
async function toResult(response) {
  const body = await response.json();
  if (response.status < 400) {
    return {
      ok: true,
      data: body
    };
  }
  if (response.status === 500) {
    notify_exports.error("An unexpected error occurred while processing your request.");
  } else if (response.status === 403) {
    notify_exports.error("You are not authorized to perform this operation.");
  }
  return {
    ok: false,
    data: body,
    error: {
      errors: body.errors
    }
  };
}
async function request(url, method, body) {
  const headers = [
    ["Accept", "application/json"],
    ["Content-Type", "application/json"]
  ];
  try {
    const response = await fetch(url, {
      method,
      headers,
      body: JSON.stringify(body),
      credentials: "same-origin"
    });
    return await toResult(response);
  } catch (err) {
    const truncatedBody = truncate(body ? JSON.stringify(body) : "<empty>", 1e3);
    throw new Error(`Failed to ${method} ${url} with body '${truncatedBody}'`);
  }
}
var http = {
  get: async (url) => {
    return await request(url, "GET");
  },
  post: async (url, body) => {
    return await request(url, "POST", body);
  },
  put: async (url, body) => {
    return await request(url, "PUT", body);
  },
  delete: async (url, body) => {
    return await request(url, "DELETE", body);
  },
  event: (category, action) => (result) => {
    if (result && result.ok) {
      analytics.event(category, action);
    }
    return result;
  }
};

// public/services/cache.ts
var set = (storage, key, value) => {
  if (storage) {
    storage.setItem(key, value);
  }
};
var get = (storage, key) => {
  if (window.localStorage) {
    return storage.getItem(key);
  }
  return null;
};
var has = (storage, key) => {
  if (storage) {
    return !!storage.getItem(key);
  }
  return false;
};
var remove = (storage, ...keys) => {
  if (storage && keys) {
    for (const key of keys) {
      storage.removeItem(key);
    }
  }
};
var cache = {
  local: {
    set: (key, value) => {
      set(window.localStorage, key, value);
    },
    get: (key) => {
      return get(window.localStorage, key);
    },
    has: (key) => {
      return has(window.localStorage, key);
    },
    remove: (...keys) => {
      remove(window.localStorage, ...keys);
    }
  },
  session: {
    set: (key, value) => {
      set(window.sessionStorage, key, value);
    },
    get: (key) => {
      return typeof window !== "undefined" && get(window.sessionStorage, key);
    },
    has: (key) => {
      return typeof window !== "undefined" && has(window.sessionStorage, key);
    },
    remove: (...keys) => {
      remove(window.sessionStorage, ...keys);
    }
  }
};

// public/services/analytics.ts
var analytics = {
  event: (eventCategory, eventAction) => {
    if (window.ga) {
      window.ga("send", "event", {
        eventCategory,
        eventAction
      });
    }
  },
  error: (err) => {
    if (window.ga) {
      window.ga("send", "exception", {
        exDescription: err ? err.stack : "<not available>",
        exFatal: false
      });
    }
  }
};

// public/services/fider.ts
var import_react2 = __toModule(require_react());
var FiderSession = class {
  constructor(data) {
    this.pProps = {};
    this.pContextID = data.contextID;
    this.pProps = data.props;
    this.pUser = data.user;
    this.pTenant = data.tenant;
  }
  get contextID() {
    return this.pContextID;
  }
  get user() {
    return this.pUser;
  }
  get tenant() {
    return this.pTenant;
  }
  get props() {
    return this.pProps;
  }
  get isAuthenticated() {
    return !!this.pUser;
  }
};
var FiderImpl = class {
  constructor() {
    this.initialize = (d) => {
      if (d) {
        this.pSettings = d.settings;
        this.pSession = new FiderSession(d);
        return this;
      }
      const el = document.getElementById("server-data");
      const data = el ? JSON.parse(el.textContent || el.innerText) : {};
      this.pSettings = data.settings;
      this.pSession = new FiderSession(data);
      return this;
    };
  }
  get session() {
    return this.pSession;
  }
  get settings() {
    return this.pSettings;
  }
  isProduction() {
    return this.pSettings.environment === "production";
  }
  isSingleHostMode() {
    return this.pSettings.mode === "single";
  }
};
var Fider = new FiderImpl();
var FiderContext = (0, import_react2.createContext)(Fider);

// public/services/utils.ts
var classSet = (input) => {
  let classes = "";
  if (input) {
    for (const key in input) {
      if (key && !!input[key]) {
        classes += ` ${key}`;
      }
    }
    return classes.trim();
  }
  return "";
};
var fileToBase64 = async (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.addEventListener("load", () => {
      const parts = reader.result.split("base64,");
      resolve(parts[1]);
    }, false);
    reader.addEventListener("error", () => {
      reject(reader.error);
    }, false);
    reader.readAsDataURL(file);
  });
};
var isCookieEnabled = () => {
  try {
    document.cookie = "cookietest=1";
    const ret = document.cookie.indexOf("cookietest=") !== -1;
    document.cookie = "cookietest=1; expires=Thu, 01-Jan-1970 00:00:01 GMT";
    return ret;
  } catch (e) {
    return false;
  }
};
var uploadedImageURL = (bkey, size) => {
  if (bkey) {
    if (size) {
      return `${Fider.settings.tenantAssetsURL}/images/${bkey}?size=${size}`;
    }
    return `${Fider.settings.tenantAssetsURL}/images/${bkey}`;
  }
  return void 0;
};
var truncate = (input, maxLength) => {
  if (input && input.length > 1e3) {
    return `${input.substr(0, maxLength)}...`;
  }
  return input;
};

// public/services/markdown.ts
var markdown_exports = {};
__export(markdown_exports, {
  full: () => full,
  simple: () => simple
});
var import_marked = __toModule(require_marked());
var import_dompurify = __toModule(require_purify_cjs());
import_marked.default.setOptions({
  headerIds: false,
  xhtml: true,
  smartLists: true,
  gfm: true,
  breaks: true
});
console.log(import_dompurify.default)
import_dompurify.default.setConfig({
  ADD_ATTR: ["target"]
});
var link = (href, title, text) => {
  const titleAttr = title ? ` title=${title}` : "";
  return `<a href="${href}"${titleAttr} rel="noopener" target="_blank">${text}</a>`;
};
var simpleRenderer = new import_marked.default.Renderer();
simpleRenderer.heading = (text, level, raw) => `<p>${raw}</p>`;
simpleRenderer.image = (href, title, text) => "";
simpleRenderer.link = link;
var fullRenderer = new import_marked.default.Renderer();
fullRenderer.link = link;
var entities = {
  "<": "&lt;",
  ">": "&gt;"
};
var encodeHTML = (s) => s.replace(/[<>]/g, (tag) => entities[tag] || tag);
var full = (input) => {
  return (0, import_marked.default)(encodeHTML(input), {renderer: fullRenderer}).trim();
};
var simple = (input) => {
  return (0, import_marked.default)(encodeHTML(input), {renderer: simpleRenderer}).trim();
};

// public/services/notify.ts
var notify_exports = {};
__export(notify_exports, {
  error: () => error,
  success: () => success
});
var toastify = () => Promise.resolve().then(() => require_toastify());
var success = (content) => {
  return toastify().then((toast) => {
    toast.success(content);
  });
};
var error = (content) => {
  return toastify().then((toast) => {
    toast.error(content);
  });
};

// public/services/querystring.ts
var querystring_exports = {};
__export(querystring_exports, {
  get: () => get2,
  getArray: () => getArray,
  getNumber: () => getNumber,
  set: () => set2,
  stringify: () => stringify
});

// public/services/navigator.ts
var navigator2 = {
  url: () => {
    return typeof window === "undefined" ? "" : window.location.href;
  },
  goHome: () => {
    window.location.href = "/";
  },
  goTo: (url) => {
    const isEqual = window.location.href === url || window.location.pathname === url;
    if (!isEqual) {
      window.location.href = url;
    }
  },
  replaceState: (path) => {
    if (history.replaceState !== void 0) {
      const newURL = Fider.settings.baseURL + path;
      window.history.replaceState({path: newURL}, "", newURL);
    }
  }
};
var navigator_default = navigator2;

// public/services/querystring.ts
var getNumber = (name) => {
  return parseInt(get2(name), 10) || void 0;
};
var set2 = (name, value) => {
  const uri = navigator_default.url();
  const re = new RegExp("([?&])" + name + "=.*?(&|$)", "i");
  if (uri.match(re)) {
    return uri.replace(re, "$1" + name + "=" + value + "$2");
  } else {
    const separator = uri.indexOf("?") !== -1 ? "&" : "?";
    return uri + separator + name + "=" + value;
  }
};
var get2 = (name) => {
  name = name.replace(/[\[\]]/g, "\\$&");
  const regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)");
  const results = regex.exec(navigator_default.url());
  if (!results || !results[2]) {
    return "";
  }
  return decodeURIComponent(results[2].replace(/\+/g, " "));
};
var getArray = (name) => {
  const qs = get2(name);
  if (qs) {
    return qs.split(",").filter((i) => i);
  }
  return [];
};
var stringify = (object) => {
  if (!object) {
    return "";
  }
  let qs = "";
  for (const key of Object.keys(object)) {
    const symbol = qs ? "&" : "?";
    const value = object[key];
    if (value instanceof Array) {
      if (value.length > 0) {
        qs += `${symbol}${key}=${value.join(",")}`;
      }
    } else if (value) {
      qs += `${symbol}${key}=${encodeURIComponent(value.toString()).replace(/%20/g, "+")}`;
    }
  }
  return qs;
};

// public/services/device.ts
var device_exports = {};
__export(device_exports, {
  isTouch: () => isTouch
});
var isTouch = () => {
  return typeof window === "undefined" ? false : "ontouchstart" in window || navigator.maxTouchPoints > 0 || navigator.msMaxTouchPoints > 0;
};

// public/services/actions/index.ts
var actions_exports = {};
__export(actions_exports, {
  addVote: () => addVote,
  assignTag: () => assignTag,
  blockUser: () => blockUser,
  changeUserEmail: () => changeUserEmail,
  changeUserRole: () => changeUserRole,
  checkAvailability: () => checkAvailability,
  completeProfile: () => completeProfile,
  createComment: () => createComment,
  createPost: () => createPost,
  createTag: () => createTag,
  createTenant: () => createTenant,
  deleteComment: () => deleteComment,
  deleteCurrentAccount: () => deleteCurrentAccount,
  deletePost: () => deletePost,
  deleteTag: () => deleteTag,
  getAllPosts: () => getAllPosts,
  getOAuthConfig: () => getOAuthConfig,
  getTotalUnreadNotifications: () => getTotalUnreadNotifications,
  listVotes: () => listVotes,
  logError: () => logError,
  markAllAsRead: () => markAllAsRead,
  regenerateAPIKey: () => regenerateAPIKey,
  removeVote: () => removeVote,
  respond: () => respond,
  saveOAuthConfig: () => saveOAuthConfig,
  searchPosts: () => searchPosts,
  sendInvites: () => sendInvites,
  sendSampleInvite: () => sendSampleInvite,
  signIn: () => signIn,
  subscribe: () => subscribe,
  unassignTag: () => unassignTag,
  unblockUser: () => unblockUser,
  unsubscribe: () => unsubscribe,
  updateComment: () => updateComment,
  updatePost: () => updatePost,
  updateTag: () => updateTag,
  updateTenantAdvancedSettings: () => updateTenantAdvancedSettings,
  updateTenantPrivacy: () => updateTenantPrivacy,
  updateTenantSettings: () => updateTenantSettings,
  updateUserSettings: () => updateUserSettings
});

// public/services/actions/user.ts
var updateUserSettings = async (request2) => {
  return await http.post("/_api/user/settings", request2);
};
var changeUserEmail = async (email) => {
  return await http.post("/_api/user/change-email", {
    email
  });
};
var deleteCurrentAccount = async () => {
  return await http.delete("/_api/user");
};
var regenerateAPIKey = async () => {
  return await http.post("/_api/user/regenerate-apikey");
};

// public/services/actions/tag.ts
var createTag = async (name, color, isPublic) => {
  return http.post(`/api/v1/tags`, {name, color, isPublic}).then(http.event("tag", "create"));
};
var updateTag = async (slug, name, color, isPublic) => {
  return http.put(`/api/v1/tags/${slug}`, {name, color, isPublic}).then(http.event("tag", "update"));
};
var deleteTag = async (slug) => {
  return http.delete(`/api/v1/tags/${slug}`).then(http.event("tag", "delete"));
};
var assignTag = async (slug, postNumber) => {
  return http.post(`/api/v1/posts/${postNumber}/tags/${slug}`).then(http.event("tag", "assign"));
};
var unassignTag = async (slug, postNumber) => {
  return http.delete(`/api/v1/posts/${postNumber}/tags/${slug}`).then(http.event("tag", "unassign"));
};

// public/services/actions/post.ts
var getAllPosts = async () => {
  return await http.get("/api/v1/posts");
};
var searchPosts = async (params) => {
  return await http.get(`/api/v1/posts${querystring_exports.stringify({
    tags: params.tags,
    query: params.query,
    view: params.view,
    limit: params.limit
  })}`);
};
var deletePost = async (postNumber, text) => {
  return http.delete(`/api/v1/posts/${postNumber}`, {
    text
  }).then(http.event("post", "delete"));
};
var addVote = async (postNumber) => {
  return http.post(`/api/v1/posts/${postNumber}/votes`).then(http.event("post", "vote"));
};
var removeVote = async (postNumber) => {
  return http.delete(`/api/v1/posts/${postNumber}/votes`).then(http.event("post", "unvote"));
};
var subscribe = async (postNumber) => {
  return http.post(`/api/v1/posts/${postNumber}/subscription`).then(http.event("post", "subscribe"));
};
var unsubscribe = async (postNumber) => {
  return http.delete(`/api/v1/posts/${postNumber}/subscription`).then(http.event("post", "unsubscribe"));
};
var listVotes = async (postNumber) => {
  return http.get(`/api/v1/posts/${postNumber}/votes`);
};
var createComment = async (postNumber, content, attachments) => {
  return http.post(`/api/v1/posts/${postNumber}/comments`, {content, attachments}).then(http.event("comment", "create"));
};
var updateComment = async (postNumber, commentID, content, attachments) => {
  return http.put(`/api/v1/posts/${postNumber}/comments/${commentID}`, {content, attachments}).then(http.event("comment", "update"));
};
var deleteComment = async (postNumber, commentID) => {
  return http.delete(`/api/v1/posts/${postNumber}/comments/${commentID}`).then(http.event("comment", "delete"));
};
var respond = async (postNumber, input) => {
  return http.put(`/api/v1/posts/${postNumber}/status`, {
    status: input.status,
    text: input.text,
    originalNumber: input.originalNumber
  }).then(http.event("post", "respond"));
};
var createPost = async (title, description, attachments) => {
  return http.post(`/api/v1/posts`, {title, description, attachments}).then(http.event("post", "create"));
};
var updatePost = async (postNumber, title, description, attachments) => {
  return http.put(`/api/v1/posts/${postNumber}`, {title, description, attachments}).then(http.event("post", "update"));
};

// public/services/actions/tenant.ts
var createTenant = async (request2) => {
  return await http.post("/_api/tenants", request2);
};
var updateTenantSettings = async (request2) => {
  return await http.post("/_api/admin/settings/general", request2);
};
var updateTenantAdvancedSettings = async (customCSS) => {
  return await http.post("/_api/admin/settings/advanced", {customCSS});
};
var updateTenantPrivacy = async (isPrivate) => {
  return await http.post("/_api/admin/settings/privacy", {
    isPrivate
  });
};
var checkAvailability = async (subdomain) => {
  return await http.get(`/_api/tenants/${subdomain}/availability`);
};
var signIn = async (email) => {
  return await http.post("/_api/signin", {
    email
  });
};
var completeProfile = async (key, name) => {
  return await http.post("/_api/signin/complete", {
    key,
    name
  });
};
var changeUserRole = async (userID, role) => {
  return await http.post(`/_api/admin/roles/${role}/users`, {
    userID
  });
};
var blockUser = async (userID) => {
  return await http.put(`/_api/admin/users/${userID}/block`);
};
var unblockUser = async (userID) => {
  return await http.delete(`/_api/admin/users/${userID}/block`);
};
var getOAuthConfig = async (provider) => {
  return await http.get(`/_api/admin/oauth/${provider}`);
};
var saveOAuthConfig = async (request2) => {
  return await http.post("/_api/admin/oauth", request2);
};

// public/services/actions/notification.ts
var getTotalUnreadNotifications = async () => {
  return http.get("/_api/notifications/unread/total").then((result) => {
    return {
      ok: result.ok,
      error: result.error,
      data: result.data ? result.data.total : 0
    };
  });
};
var markAllAsRead = async () => {
  return await http.post("/_api/notifications/read-all");
};

// public/services/actions/invite.ts
var sendInvites = async (subject, message, recipients) => {
  return http.post("/api/v1/invitations/send", {subject, message, recipients}).then(http.event("invite", "send"));
};
var sendSampleInvite = async (subject, message) => {
  return http.post("/api/v1/invitations/sample", {subject, message}).then(http.event("invite", "sample"));
};

// public/services/actions/infra.ts
var ignoreErrors = [
  "http://gj.track.uc.cn/collect",
  "null is not an object (evaluating 'c.sheet.insertRule')",
  "Refused to evaluate a string as JavaScript because 'unsafe-eval'",
  "vid_mate_check is not defined",
  "SecurityError: Failed to read the 'cssRules' property from 'CSSStyleSheet': Cannot access rules"
];
var logError = async (message, err) => {
  for (const pattern of ignoreErrors) {
    if (message.indexOf(pattern) >= 0) {
      return;
    }
  }
  const data = {
    url: navigator_default.url(),
    stack: err ? err.stack : "<not available>"
  };
  try {
    analytics.error(err);
    return await http.post("/_api/log-error", {message, data});
  } catch (err2) {
  }
};

// public/hooks/use-fider.ts
var useFider = () => (0, import_react3.useContext)(FiderContext);

// public/pages/Error/Error.page.tsx
var ErrorPage = (props) => {
  const fider = useFider();
  return /* @__PURE__ */ import_react4.default.createElement("div", {
    id: "p-error",
    className: "container failure-page"
  }, /* @__PURE__ */ import_react4.default.createElement(TenantLogo, {
    size: 100,
    useFiderIfEmpty: true
  }), /* @__PURE__ */ import_react4.default.createElement("h1", null, "Shoot! Well, this is unexpected\u2026"), /* @__PURE__ */ import_react4.default.createElement("p", null, "An error has occurred and we're working to fix the problem!"), fider.settings && /* @__PURE__ */ import_react4.default.createElement("span", null, "Take me back to ", /* @__PURE__ */ import_react4.default.createElement("a", {
    href: fider.settings.baseURL
  }, fider.settings.baseURL), " home page."), props.showDetails && /* @__PURE__ */ import_react4.default.createElement("pre", {
    className: "error"
  }, props.error.toString(), props.errorInfo.componentStack));
};

// public/components/ErrorBoundary.tsx
var ErrorBoundary = class extends import_react5.default.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: void 0,
      errorInfo: void 0
    };
  }
  componentDidCatch(error2, errorInfo) {
    const onError = this.props.onError;
    if (onError) {
      onError(error2);
    }
    this.setState({
      error: error2,
      errorInfo
    });
  }
  render() {
    const {error: error2, errorInfo} = this.state;
    if (error2 && errorInfo) {
      return /* @__PURE__ */ import_react5.default.createElement(FiderContext.Consumer, null, (fider) => /* @__PURE__ */ import_react5.default.createElement(ErrorPage, {
        error: error2,
        errorInfo,
        showDetails: !fider.isProduction()
      }));
    } else {
      return this.props.children;
    }
  }
};

// public/components/ShowPostResponse.tsx
var import_react39 = __toModule(require_react());

// public/components/common/Button.tsx
var import_react6 = __toModule(require_react());
var ButtonClickEvent = class {
  constructor() {
    this.shouldEnable = true;
  }
  preventEnable() {
    this.shouldEnable = false;
  }
  canEnable() {
    return this.shouldEnable;
  }
};
var Button = class extends import_react6.default.Component {
  constructor(props) {
    super(props);
    this.unmounted = false;
    this.click = async (e) => {
      if (e) {
        e.preventDefault();
        e.stopPropagation();
      }
      if (this.state.clicked) {
        return;
      }
      const event = new ButtonClickEvent();
      this.setState({clicked: true});
      if (this.props.onClick) {
        await this.props.onClick(event);
        if (!this.unmounted && event.canEnable()) {
          this.setState({clicked: false});
        }
      }
    };
    this.state = {
      clicked: false
    };
  }
  componentWillUnmount() {
    this.unmounted = true;
  }
  render() {
    const className = classSet({
      "c-button": true,
      "m-fluid": this.props.fluid,
      [`m-${this.props.size}`]: this.props.size,
      [`m-${this.props.color}`]: this.props.color,
      "m-loading": this.state.clicked,
      "m-disabled": this.state.clicked || this.props.disabled,
      [this.props.className]: this.props.className
    });
    if (this.props.href) {
      return /* @__PURE__ */ import_react6.default.createElement("a", {
        href: this.props.href,
        rel: this.props.rel,
        className
      }, this.props.children);
    } else if (this.props.onClick) {
      return /* @__PURE__ */ import_react6.default.createElement("button", {
        type: this.props.type,
        className,
        onClick: this.click
      }, this.props.children);
    } else {
      return /* @__PURE__ */ import_react6.default.createElement("button", {
        type: this.props.type,
        className
      }, this.props.children);
    }
  }
};
Button.defaultProps = {
  size: "small",
  fluid: false,
  color: "default",
  type: "button"
};

// public/components/common/form/Form.tsx
var import_react7 = __toModule(require_react());
var ValidationContext = import_react7.default.createContext({});
var Form = (props) => {
  const className = classSet({
    "c-form": true,
    [props.className]: props.className,
    [`m-${props.size}`]: props.size
  });
  return /* @__PURE__ */ import_react7.default.createElement("form", {
    autoComplete: "off",
    className
  }, /* @__PURE__ */ import_react7.default.createElement(DisplayError, {
    error: props.error
  }), /* @__PURE__ */ import_react7.default.createElement(ValidationContext.Provider, {
    value: {error: props.error}
  }, props.children));
};

// public/components/common/form/Input.tsx
var import_react9 = __toModule(require_react());

// public/components/common/form/DisplayError.tsx
var import_react8 = __toModule(require_react());
var arrayToTag = (items) => {
  return items.map((m) => /* @__PURE__ */ import_react8.default.createElement("li", {
    key: m
  }, m));
};
var hasError = (field, error2) => {
  if (field && error2 && error2.errors) {
    for (const err of error2.errors) {
      if (err.field === field) {
        return true;
      }
    }
  }
  return false;
};
var DisplayError = (props) => {
  if (!props.error || !props.error.errors) {
    return null;
  }
  const dict = props.error.errors.reduce((result, err) => {
    result[err.field || ""] = result[err.field || ""] || [];
    result[err.field || ""].push(err.message);
    return result;
  }, {});
  let items = [];
  if (dict[""] && !props.fields) {
    items = arrayToTag(dict[""]);
  } else if (props.fields) {
    for (const field of props.fields || Object.keys(dict)) {
      if (dict.hasOwnProperty(field)) {
        const tags = arrayToTag(dict[field]);
        tags.forEach((t) => items.push(t));
      }
    }
  }
  return items.length > 0 ? /* @__PURE__ */ import_react8.default.createElement("div", {
    className: `c-form-field-error`
  }, /* @__PURE__ */ import_react8.default.createElement("ul", null, items)) : null;
};

// public/components/common/form/Input.tsx
var Input = (props) => {
  const onChange = (e) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value);
    }
  };
  const suffix = typeof props.suffix === "string" ? /* @__PURE__ */ import_react9.default.createElement("span", {
    className: "c-form-input-suffix"
  }, props.suffix) : props.suffix;
  const icon = !!props.icon ? import_react9.default.createElement(props.icon, {
    onClick: props.onIconClick,
    className: classSet({link: !!props.onIconClick})
  }) : void 0;
  return /* @__PURE__ */ import_react9.default.createElement(ValidationContext.Consumer, null, (ctx) => /* @__PURE__ */ import_react9.default.createElement("div", {
    className: classSet({
      "c-form-field": true,
      "m-suffix": props.suffix,
      "m-error": hasError(props.field, ctx.error),
      "m-icon": !!props.icon,
      [`${props.className}`]: props.className
    })
  }, !!props.label && /* @__PURE__ */ import_react9.default.createElement("label", {
    htmlFor: `input-${props.field}`
  }, props.label, props.afterLabel), /* @__PURE__ */ import_react9.default.createElement("div", {
    className: "c-form-field-wrapper"
  }, /* @__PURE__ */ import_react9.default.createElement("input", {
    id: `input-${props.field}`,
    type: "text",
    autoComplete: props.autoComplete,
    tabIndex: props.noTabFocus ? -1 : void 0,
    ref: props.inputRef,
    autoFocus: props.autoFocus,
    onFocus: props.onFocus,
    maxLength: props.maxLength,
    disabled: props.disabled,
    value: props.value,
    placeholder: props.placeholder,
    onChange
  }), icon, suffix), /* @__PURE__ */ import_react9.default.createElement(DisplayError, {
    fields: [props.field],
    error: ctx.error
  }), props.children));
};

// public/components/common/form/ImageUploader.tsx
var import_react10 = __toModule(require_react());

// node_modules/react-icons/fa/index.esm.js
var import_lib = __toModule(require_cjs());
function FaBan(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 512 512"}, child: [{tag: "path", attr: {d: "M256 8C119.034 8 8 119.033 8 256s111.034 248 248 248 248-111.034 248-248S392.967 8 256 8zm130.108 117.892c65.448 65.448 70 165.481 20.677 235.637L150.47 105.216c70.204-49.356 170.226-44.735 235.638 20.676zM125.892 386.108c-65.448-65.448-70-165.481-20.677-235.637L361.53 406.784c-70.203 49.356-170.226 44.736-235.638-20.676z"}}]})(props);
}
function FaCaretUp(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 320 512"}, child: [{tag: "path", attr: {d: "M288.662 352H31.338c-17.818 0-26.741-21.543-14.142-34.142l128.662-128.662c7.81-7.81 20.474-7.81 28.284 0l128.662 128.662c12.6 12.599 3.676 34.142-14.142 34.142z"}}]})(props);
}
function FaCheck(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 512 512"}, child: [{tag: "path", attr: {d: "M173.898 439.404l-166.4-166.4c-9.997-9.997-9.997-26.206 0-36.204l36.203-36.204c9.997-9.998 26.207-9.998 36.204 0L192 312.69 432.095 72.596c9.997-9.997 26.207-9.997 36.204 0l36.203 36.204c9.997 9.997 9.997 26.206 0 36.204l-294.4 294.401c-9.998 9.997-26.207 9.997-36.204-.001z"}}]})(props);
}
function FaExclamationTriangle(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 576 512"}, child: [{tag: "path", attr: {d: "M569.517 440.013C587.975 472.007 564.806 512 527.94 512H48.054c-36.937 0-59.999-40.055-41.577-71.987L246.423 23.985c18.467-32.009 64.72-31.951 83.154 0l239.94 416.028zM288 354c-25.405 0-46 20.595-46 46s20.595 46 46 46 46-20.595 46-46-20.595-46-46-46zm-43.673-165.346l7.418 136c.347 6.364 5.609 11.346 11.982 11.346h48.546c6.373 0 11.635-4.982 11.982-11.346l7.418-136c.375-6.874-5.098-12.654-11.982-12.654h-63.383c-6.884 0-12.356 5.78-11.981 12.654z"}}]})(props);
}
function FaLock(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 448 512"}, child: [{tag: "path", attr: {d: "M400 224h-24v-72C376 68.2 307.8 0 224 0S72 68.2 72 152v72H48c-26.5 0-48 21.5-48 48v192c0 26.5 21.5 48 48 48h352c26.5 0 48-21.5 48-48V272c0-26.5-21.5-48-48-48zm-104 0H152v-72c0-39.7 32.3-72 72-72s72 32.3 72 72v72z"}}]})(props);
}
function FaSearch(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 512 512"}, child: [{tag: "path", attr: {d: "M505 442.7L405.3 343c-4.5-4.5-10.6-7-17-7H372c27.6-35.3 44-79.7 44-128C416 93.1 322.9 0 208 0S0 93.1 0 208s93.1 208 208 208c48.3 0 92.7-16.4 128-44v16.3c0 6.4 2.5 12.5 7 17l99.7 99.7c9.4 9.4 24.6 9.4 33.9 0l28.3-28.3c9.4-9.4 9.4-24.6.1-34zM208 336c-70.7 0-128-57.2-128-128 0-70.7 57.2-128 128-128 70.7 0 128 57.2 128 128 0 70.7-57.2 128-128 128z"}}]})(props);
}
function FaTimes(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 352 512"}, child: [{tag: "path", attr: {d: "M242.72 256l100.07-100.07c12.28-12.28 12.28-32.19 0-44.48l-22.24-22.24c-12.28-12.28-32.19-12.28-44.48 0L176 189.28 75.93 89.21c-12.28-12.28-32.19-12.28-44.48 0L9.21 111.45c-12.28 12.28-12.28 32.19 0 44.48L109.28 256 9.21 356.07c-12.28 12.28-12.28 32.19 0 44.48l22.24 22.24c12.28 12.28 32.2 12.28 44.48 0L176 322.72l100.07 100.07c12.28 12.28 32.2 12.28 44.48 0l22.24-22.24c12.28-12.28 12.28-32.19 0-44.48L242.72 256z"}}]})(props);
}
function FaRegCheckCircle(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 512 512"}, child: [{tag: "path", attr: {d: "M256 8C119.033 8 8 119.033 8 256s111.033 248 248 248 248-111.033 248-248S392.967 8 256 8zm0 48c110.532 0 200 89.451 200 200 0 110.532-89.451 200-200 200-110.532 0-200-89.451-200-200 0-110.532 89.451-200 200-200m140.204 130.267l-22.536-22.718c-4.667-4.705-12.265-4.736-16.97-.068L215.346 303.697l-59.792-60.277c-4.667-4.705-12.265-4.736-16.97-.069l-22.719 22.536c-4.705 4.667-4.736 12.265-.068 16.971l90.781 91.516c4.667 4.705 12.265 4.736 16.97.068l172.589-171.204c4.704-4.668 4.734-12.266.067-16.971z"}}]})(props);
}
function FaRegComments(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 576 512"}, child: [{tag: "path", attr: {d: "M532 386.2c27.5-27.1 44-61.1 44-98.2 0-80-76.5-146.1-176.2-157.9C368.3 72.5 294.3 32 208 32 93.1 32 0 103.6 0 192c0 37 16.5 71 44 98.2-15.3 30.7-37.3 54.5-37.7 54.9-6.3 6.7-8.1 16.5-4.4 25 3.6 8.5 12 14 21.2 14 53.5 0 96.7-20.2 125.2-38.8 9.2 2.1 18.7 3.7 28.4 4.9C208.1 407.6 281.8 448 368 448c20.8 0 40.8-2.4 59.8-6.8C456.3 459.7 499.4 480 553 480c9.2 0 17.5-5.5 21.2-14 3.6-8.5 1.9-18.3-4.4-25-.4-.3-22.5-24.1-37.8-54.8zm-392.8-92.3L122.1 305c-14.1 9.1-28.5 16.3-43.1 21.4 2.7-4.7 5.4-9.7 8-14.8l15.5-31.1L77.7 256C64.2 242.6 48 220.7 48 192c0-60.7 73.3-112 160-112s160 51.3 160 112-73.3 112-160 112c-16.5 0-33-1.9-49-5.6l-19.8-4.5zM498.3 352l-24.7 24.4 15.5 31.1c2.6 5.1 5.3 10.1 8 14.8-14.6-5.1-29-12.3-43.1-21.4l-17.1-11.1-19.9 4.6c-16 3.7-32.5 5.6-49 5.6-54 0-102.2-20.1-131.3-49.7C338 339.5 416 272.9 416 192c0-3.4-.4-6.7-.7-10C479.7 196.5 528 238.8 528 288c0 28.7-16.2 50.6-29.7 64z"}}]})(props);
}
function FaRegImage(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 512 512"}, child: [{tag: "path", attr: {d: "M464 64H48C21.49 64 0 85.49 0 112v288c0 26.51 21.49 48 48 48h416c26.51 0 48-21.49 48-48V112c0-26.51-21.49-48-48-48zm-6 336H54a6 6 0 0 1-6-6V118a6 6 0 0 1 6-6h404a6 6 0 0 1 6 6v276a6 6 0 0 1-6 6zM128 152c-22.091 0-40 17.909-40 40s17.909 40 40 40 40-17.909 40-40-17.909-40-40-40zM96 352h320v-80l-87.515-87.515c-4.686-4.686-12.284-4.686-16.971 0L192 304l-39.515-39.515c-4.686-4.686-12.284-4.686-16.971 0L96 304v48z"}}]})(props);
}
function FaRegLightbulb(props) {
  return (0, import_lib.GenIcon)({tag: "svg", attr: {viewBox: "0 0 352 512"}, child: [{tag: "path", attr: {d: "M176 80c-52.94 0-96 43.06-96 96 0 8.84 7.16 16 16 16s16-7.16 16-16c0-35.3 28.72-64 64-64 8.84 0 16-7.16 16-16s-7.16-16-16-16zM96.06 459.17c0 3.15.93 6.22 2.68 8.84l24.51 36.84c2.97 4.46 7.97 7.14 13.32 7.14h78.85c5.36 0 10.36-2.68 13.32-7.14l24.51-36.84c1.74-2.62 2.67-5.7 2.68-8.84l.05-43.18H96.02l.04 43.18zM176 0C73.72 0 0 82.97 0 176c0 44.37 16.45 84.85 43.56 115.78 16.64 18.99 42.74 58.8 52.42 92.16v.06h48v-.12c-.01-4.77-.72-9.51-2.15-14.07-5.59-17.81-22.82-64.77-62.17-109.67-20.54-23.43-31.52-53.15-31.61-84.14-.2-73.64 59.67-128 127.95-128 70.58 0 128 57.42 128 128 0 30.97-11.24 60.85-31.65 84.14-39.11 44.61-56.42 91.47-62.1 109.46a47.507 47.507 0 0 0-2.22 14.3v.1h48v-.05c9.68-33.37 35.78-73.18 52.42-92.16C335.55 260.85 352 220.37 352 176 352 78.8 273.2 0 176 0z"}}]})(props);
}

// public/components/common/form/ImageUploader.tsx
var hardFileSizeLimit = 5 * 1024 * 1024;
var ImageUploader = class extends import_react10.default.Component {
  constructor(props) {
    super(props);
    this.fileChanged = async (e) => {
      if (e.target.files && e.target.files[0]) {
        const file = e.target.files[0];
        if (file.size > hardFileSizeLimit) {
          alert("The image size must be smaller than 5MB.");
          return;
        }
        const base64 = await fileToBase64(file);
        this.setState({
          bkey: this.props.bkey,
          upload: {
            fileName: file.name,
            content: base64,
            contentType: file.type
          },
          remove: false,
          previewURL: `data:${file.type};base64,${base64}`
        }, () => {
          this.props.onChange(this.state, this.props.instanceID, this.state.previewURL);
        });
      }
    };
    this.removeFile = async (e) => {
      if (this.fileSelector) {
        this.fileSelector.value = "";
      }
      this.setState({
        bkey: this.props.bkey,
        remove: true,
        upload: void 0,
        previewURL: void 0
      }, () => {
        this.props.onChange({
          bkey: this.state.bkey,
          remove: this.state.remove,
          upload: this.state.upload
        }, this.props.instanceID, this.state.previewURL);
      });
    };
    this.selectFile = async (e) => {
      if (this.fileSelector) {
        this.fileSelector.click();
      }
    };
    this.openModal = () => {
      this.setState({showModal: true});
    };
    this.closeModal = async () => {
      this.setState({showModal: false});
    };
    this.state = {
      upload: void 0,
      remove: false,
      showModal: false,
      previewURL: uploadedImageURL(this.props.bkey, this.props.previewMaxWidth)
    };
  }
  modal() {
    return /* @__PURE__ */ import_react10.default.createElement(Modal.Window, {
      className: "c-image-viewer-modal",
      isOpen: this.state.showModal,
      onClose: this.closeModal,
      center: false,
      size: "fluid"
    }, /* @__PURE__ */ import_react10.default.createElement(Modal.Content, null, this.props.bkey ? /* @__PURE__ */ import_react10.default.createElement("img", {
      src: uploadedImageURL(this.props.bkey)
    }) : /* @__PURE__ */ import_react10.default.createElement("img", {
      src: this.state.previewURL
    })), /* @__PURE__ */ import_react10.default.createElement(Modal.Footer, null, /* @__PURE__ */ import_react10.default.createElement(Button, {
      color: "cancel",
      onClick: this.closeModal
    }, "Close")));
  }
  render() {
    const isUploading = !!this.state.upload;
    const hasFile = !this.state.remove && this.props.bkey || isUploading;
    const imgStyles = {
      maxWidth: `${this.props.previewMaxWidth}px`
    };
    return /* @__PURE__ */ import_react10.default.createElement(ValidationContext.Consumer, null, (ctx) => /* @__PURE__ */ import_react10.default.createElement("div", {
      className: classSet({
        "c-form-field": true,
        "c-image-upload": true,
        "m-error": hasError(this.props.field, ctx.error)
      })
    }, this.modal(), /* @__PURE__ */ import_react10.default.createElement("label", {
      htmlFor: `input-${this.props.field}`
    }, this.props.label), hasFile && /* @__PURE__ */ import_react10.default.createElement("div", {
      className: "preview"
    }, /* @__PURE__ */ import_react10.default.createElement("img", {
      onClick: this.openModal,
      src: this.state.previewURL,
      style: imgStyles
    }), !this.props.disabled && /* @__PURE__ */ import_react10.default.createElement(Button, {
      onClick: this.removeFile,
      color: "danger"
    }, "X")), /* @__PURE__ */ import_react10.default.createElement("input", {
      ref: (e) => this.fileSelector = e,
      type: "file",
      onChange: this.fileChanged,
      accept: "image/*"
    }), /* @__PURE__ */ import_react10.default.createElement(DisplayError, {
      fields: [this.props.field],
      error: ctx.error
    }), !hasFile && /* @__PURE__ */ import_react10.default.createElement("div", {
      className: "c-form-field-wrapper"
    }, /* @__PURE__ */ import_react10.default.createElement(Button, {
      onClick: this.selectFile,
      disabled: this.props.disabled
    }, /* @__PURE__ */ import_react10.default.createElement(FaRegImage, null))), this.props.children));
  }
};

// public/components/common/form/MultiImageUploader.tsx
var import_react11 = __toModule(require_react());
var MultiImageUploader = class extends import_react11.default.Component {
  constructor(props) {
    super(props);
    this.imageUploaded = (upload, instanceID) => {
      const instances = {...this.state.instances};
      const removed = [...this.state.removed];
      let count = this.state.count;
      if (upload.remove) {
        if (upload.bkey) {
          removed.push(upload);
        }
        delete instances[instanceID];
        if (--count === this.props.maxUploads) {
          this.addNewElement(instances);
        }
      } else {
        instances[instanceID].upload = upload;
        if (count++ <= this.props.maxUploads) {
          this.addNewElement(instances);
        }
      }
      this.setState({instances, count, removed}, this.triggerOnChange);
    };
    let count = 1;
    const instances = {};
    if (props.bkeys) {
      for (const bkey of props.bkeys) {
        count++;
        this.addNewElement(instances, bkey);
      }
    }
    if (count <= this.props.maxUploads) {
      count++;
      this.addNewElement(instances);
    }
    this.state = {instances, count, removed: []};
  }
  triggerOnChange() {
    if (this.props.onChange) {
      const uploads = Object.keys(this.state.instances).map((k) => this.state.instances[k].upload).concat(this.state.removed).filter((x) => !!x);
      this.props.onChange(uploads);
    }
  }
  addNewElement(instances, bkey) {
    const id = btoa(Math.random().toString());
    instances[id] = {
      element: /* @__PURE__ */ import_react11.default.createElement(ImageUploader, {
        key: id,
        bkey,
        instanceID: id,
        field: "attachment",
        previewMaxWidth: this.props.previewMaxWidth,
        onChange: this.imageUploaded
      })
    };
  }
  render() {
    const elements = Object.keys(this.state.instances).map((k) => this.state.instances[k].element);
    return /* @__PURE__ */ import_react11.default.createElement(ValidationContext.Consumer, null, (ctx) => /* @__PURE__ */ import_react11.default.createElement("div", {
      className: classSet({
        "c-form-field": true,
        "c-multi-image-uploader": true,
        "m-error": hasError(this.props.field, ctx.error)
      })
    }, /* @__PURE__ */ import_react11.default.createElement("div", {
      className: "c-multi-image-uploader-instances"
    }, elements), /* @__PURE__ */ import_react11.default.createElement(DisplayError, {
      fields: [this.props.field],
      error: ctx.error
    })));
  }
};

// public/components/common/form/TextArea.tsx
var import_react12 = __toModule(require_react());
var import_react_textarea_autosize = __toModule(require_react_textarea_autosize_cjs());
var TextArea = (props) => {
  const onChange = (e) => {
    if (props.onChange) {
      props.onChange(e.currentTarget.value);
    }
  };
  return /* @__PURE__ */ import_react12.default.createElement(ValidationContext.Consumer, null, (ctx) => /* @__PURE__ */ import_react12.default.createElement(import_react12.default.Fragment, null, /* @__PURE__ */ import_react12.default.createElement("div", {
    className: classSet({
      "c-form-field": true,
      "m-error": hasError(props.field, ctx.error)
    })
  }, !!props.label && /* @__PURE__ */ import_react12.default.createElement("label", {
    htmlFor: `input-${props.field}`
  }, props.label), /* @__PURE__ */ import_react12.default.createElement("div", {
    className: "c-form-field-wrapper"
  }, /* @__PURE__ */ import_react12.default.createElement(import_react_textarea_autosize.default, {
    id: `input-${props.field}`,
    disabled: props.disabled,
    onChange,
    value: props.value,
    minRows: props.minRows || 3,
    placeholder: props.placeholder,
    ref: props.inputRef,
    onFocus: props.onFocus
  })), /* @__PURE__ */ import_react12.default.createElement(DisplayError, {
    fields: [props.field],
    error: ctx.error
  }), props.children)));
};

// public/components/common/form/RadioButton.tsx
var import_react13 = __toModule(require_react());
var RadioButton = class extends import_react13.default.Component {
  constructor(props) {
    super(props);
    this.onChange = (selected) => {
      this.setState({selected}, () => {
        if (this.props.onSelect) {
          this.props.onSelect(this.state.selected);
        }
      });
    };
    this.state = {
      selected: props.defaultOption || props.options[0]
    };
  }
  render() {
    const inputs = this.props.options.map((option) => {
      return /* @__PURE__ */ import_react13.default.createElement("div", {
        key: option.value,
        className: "c-form-radio-option"
      }, /* @__PURE__ */ import_react13.default.createElement("input", {
        id: `visibility-${option.value}`,
        type: "radio",
        name: `input-${this.props.field}`,
        checked: this.state.selected === option,
        onChange: this.onChange.bind(this, option)
      }), /* @__PURE__ */ import_react13.default.createElement("label", {
        htmlFor: `visibility-${option.value}`
      }, option.label));
    });
    return /* @__PURE__ */ import_react13.default.createElement("div", {
      className: "c-form-field"
    }, /* @__PURE__ */ import_react13.default.createElement("label", {
      htmlFor: `input-${this.props.field}`
    }, this.props.label), inputs);
  }
};

// public/components/common/form/Select.tsx
var import_react14 = __toModule(require_react());
var Select = class extends import_react14.default.Component {
  constructor(props) {
    super(props);
    this.onChange = (e) => {
      let selected;
      if (e.currentTarget.value) {
        const options = this.props.options.filter((o) => o.value === e.currentTarget.value);
        if (options && options.length > 0) {
          selected = options[0];
        }
      }
      this.setState({selected}, () => {
        if (this.props.onChange) {
          this.props.onChange(this.state.selected);
        }
      });
    };
    this.state = {
      selected: this.getOption(props.defaultValue)
    };
  }
  getOption(value) {
    if (value && this.props.options) {
      const filtered = this.props.options.filter((x) => x.value === value);
      if (filtered && filtered.length > 0) {
        return filtered[0];
      }
    }
  }
  render() {
    const options = this.props.options.map((option) => {
      return /* @__PURE__ */ import_react14.default.createElement("option", {
        key: option.value,
        value: option.value
      }, option.label);
    });
    return /* @__PURE__ */ import_react14.default.createElement(ValidationContext.Consumer, null, (ctx) => /* @__PURE__ */ import_react14.default.createElement(import_react14.default.Fragment, null, /* @__PURE__ */ import_react14.default.createElement("div", {
      className: classSet({
        "c-form-field": true,
        "m-error": hasError(this.props.field, ctx.error)
      })
    }, !!this.props.label && /* @__PURE__ */ import_react14.default.createElement("label", {
      htmlFor: `input-${this.props.field}`
    }, this.props.label), /* @__PURE__ */ import_react14.default.createElement("div", {
      className: "c-form-field-wrapper"
    }, /* @__PURE__ */ import_react14.default.createElement("select", {
      id: `input-${this.props.field}`,
      defaultValue: this.props.defaultValue,
      onChange: this.onChange
    }, options)), /* @__PURE__ */ import_react14.default.createElement(DisplayError, {
      fields: [this.props.field],
      error: ctx.error
    }), this.props.children)));
  }
};

// public/components/common/form/Field.tsx
var import_react15 = __toModule(require_react());
var Field = (props) => {
  const fields = props.field ? [props.field] : void 0;
  return /* @__PURE__ */ import_react15.default.createElement(ValidationContext.Consumer, null, (ctx) => /* @__PURE__ */ import_react15.default.createElement("div", {
    className: classSet({
      "c-form-field": true,
      "m-error": hasError(props.field, ctx.error),
      [props.className]: props.className
    })
  }, !!props.label && /* @__PURE__ */ import_react15.default.createElement("label", null, props.label, props.afterLabel), props.children, /* @__PURE__ */ import_react15.default.createElement(DisplayError, {
    fields,
    error: ctx.error
  })));
};

// public/components/common/form/Checkbox.tsx
var import_react16 = __toModule(require_react());

// public/components/common/form/ImageViewer.tsx
var import_react17 = __toModule(require_react());
var ImageViewer = class extends import_react17.default.Component {
  constructor(props) {
    super(props);
    this.openModal = () => {
      if (this.state.loadedThumbnail) {
        this.setState({showModal: true});
      }
    };
    this.closeModal = async () => {
      this.setState({showModal: false});
    };
    this.onThumbnailLoad = () => {
      this.setState({loadedThumbnail: true});
    };
    this.onPreviewLoad = () => {
      this.setState({loadedPreview: true});
    };
    this.state = {
      showModal: false,
      loadedThumbnail: false,
      loadedPreview: false
    };
  }
  modal() {
    return /* @__PURE__ */ import_react17.default.createElement(Modal.Window, {
      className: "c-image-viewer-modal",
      isOpen: this.state.showModal,
      onClose: this.closeModal,
      center: false,
      size: "fluid"
    }, /* @__PURE__ */ import_react17.default.createElement(Modal.Content, null, !this.state.loadedPreview && /* @__PURE__ */ import_react17.default.createElement(Loader, null), /* @__PURE__ */ import_react17.default.createElement("img", {
      onLoad: this.onPreviewLoad,
      src: uploadedImageURL(this.props.bkey, 1500)
    })), /* @__PURE__ */ import_react17.default.createElement(Modal.Footer, null, /* @__PURE__ */ import_react17.default.createElement(Button, {
      color: "cancel",
      onClick: this.closeModal
    }, "Close")));
  }
  render() {
    const previewURL = uploadedImageURL(this.props.bkey, 200);
    return /* @__PURE__ */ import_react17.default.createElement("div", {
      className: "c-image-viewer"
    }, this.modal(), !this.state.loadedThumbnail && /* @__PURE__ */ import_react17.default.createElement(Loader, null), /* @__PURE__ */ import_react17.default.createElement("img", {
      onClick: this.openModal,
      onLoad: this.onThumbnailLoad,
      src: previewURL
    }));
  }
};

// public/components/common/MultiLineText.tsx
var import_react18 = __toModule(require_react());
var MultiLineText = (props) => {
  if (!props.text) {
    return /* @__PURE__ */ import_react18.default.createElement("p", null);
  }
  const func = props.style === "full" ? markdown_exports.full : markdown_exports.simple;
  return /* @__PURE__ */ import_react18.default.createElement("div", {
    className: `markdown-body ${props.className || ""}`,
    dangerouslySetInnerHTML: {__html: func(props.text)}
  });
};

// public/components/common/EnvironmentInfo.tsx
var import_react19 = __toModule(require_react());

// public/components/common/Avatar.tsx
var import_react20 = __toModule(require_react());
var Avatar = (props) => {
  const size = props.size || "normal";
  const className = classSet({
    "c-avatar": true,
    [`m-${size}`]: true,
    "m-staff": props.user.role && isCollaborator(props.user.role)
  });
  return /* @__PURE__ */ import_react20.default.createElement("img", {
    className,
    title: props.user.name,
    src: `${props.user.avatarURL}?size=50`
  });
};

// public/components/common/Message.tsx
var import_react21 = __toModule(require_react());
var Message = (props) => {
  const className = classSet({
    "c-message": true,
    [`m-${props.type}`]: true
  });
  const icon = props.type === "error" ? /* @__PURE__ */ import_react21.default.createElement(FaBan, null) : props.type === "warning" ? /* @__PURE__ */ import_react21.default.createElement(FaExclamationTriangle, null) : /* @__PURE__ */ import_react21.default.createElement(FaRegCheckCircle, null);
  return /* @__PURE__ */ import_react21.default.createElement("p", {
    className
  }, props.showIcon === true && icon, /* @__PURE__ */ import_react21.default.createElement("span", null, props.children));
};

// public/components/common/Hint.tsx
var import_react22 = __toModule(require_react());
var Hint = (props) => {
  const cacheKey = props.permanentCloseKey ? `Hint-Closed-${props.permanentCloseKey}` : void 0;
  const [isClosed, setIsClosed] = (0, import_react22.useState)(cacheKey ? cache.local.has(cacheKey) : false);
  const close = () => {
    if (cacheKey) {
      cache.local.set(cacheKey, "true");
    }
    setIsClosed(true);
  };
  if (props.condition === false || isClosed) {
    return null;
  }
  return /* @__PURE__ */ import_react22.default.createElement("p", {
    className: "c-hint"
  }, /* @__PURE__ */ import_react22.default.createElement("strong", null, "HINT:"), " ", props.children, cacheKey && /* @__PURE__ */ import_react22.default.createElement(FaTimes, {
    onClick: close,
    className: "close"
  }));
};

// public/components/common/Footer.tsx
var import_react23 = __toModule(require_react());

// public/components/common/Header.tsx
var import_react24 = __toModule(require_react());

// public/components/common/Heading.tsx
var import_react25 = __toModule(require_react());
var Header = (props) => import_react25.default.createElement(`h${props.level}`, {className: props.className}, props.children);
var Heading = (props) => {
  const size = props.size || "normal";
  const level = size === "normal" ? 2 : 3;
  const className = classSet({
    "c-heading": true,
    "m-dividing": props.dividing || false,
    [`m-${size}`]: true,
    [`${props.className}`]: props.className
  });
  const iconClassName = classSet({
    "c-heading-icon": true,
    circular: level <= 2
  });
  const icon = props.icon && /* @__PURE__ */ import_react25.default.createElement("div", {
    className: iconClassName
  }, import_react25.default.createElement(props.icon));
  return /* @__PURE__ */ import_react25.default.createElement(Header, {
    level,
    className
  }, icon, /* @__PURE__ */ import_react25.default.createElement("div", {
    className: "c-heading-content"
  }, props.title, /* @__PURE__ */ import_react25.default.createElement("div", {
    className: "c-heading-subtitle"
  }, props.subtitle)));
};

// public/components/common/Legal.tsx
var import_react26 = __toModule(require_react());
var TermsOfService = () => {
  const fider = useFider();
  if (fider.settings.hasLegal) {
    return /* @__PURE__ */ import_react26.default.createElement("a", {
      href: "/terms",
      target: "_blank"
    }, "Terms of Service");
  }
  return null;
};
var PrivacyPolicy = () => {
  const fider = useFider();
  if (fider.settings.hasLegal) {
    return /* @__PURE__ */ import_react26.default.createElement("a", {
      href: "/privacy",
      target: "_blank"
    }, "Privacy Policy");
  }
  return null;
};
var LegalNotice = () => {
  const fider = useFider();
  if (fider.settings.hasLegal) {
    return /* @__PURE__ */ import_react26.default.createElement("p", {
      className: "info"
    }, "By signing in, you agree to the ", /* @__PURE__ */ import_react26.default.createElement(PrivacyPolicy, null), " and ", /* @__PURE__ */ import_react26.default.createElement(TermsOfService, null), ".");
  }
  return null;
};
var LegalFooter = () => {
  const fider = useFider();
  if (fider.settings.hasLegal) {
    return /* @__PURE__ */ import_react26.default.createElement(Modal.Footer, {
      align: "center"
    }, /* @__PURE__ */ import_react26.default.createElement(LegalNotice, null));
  }
  return null;
};

// public/components/common/SocialSignInButton.tsx
var import_react27 = __toModule(require_react());
var SocialSignInButton = (props) => {
  const redirectTo = props.redirectTo || location.href;
  const href = props.option.url ? `${props.option.url}?redirect=${redirectTo}` : void 0;
  const className = classSet({
    "m-social": true,
    [`m-${props.option.provider}`]: props.option.provider
  });
  return /* @__PURE__ */ import_react27.default.createElement(Button, {
    href,
    rel: "nofollow",
    fluid: true,
    className
  }, props.option.logoURL ? /* @__PURE__ */ import_react27.default.createElement("img", {
    src: props.option.logoURL
  }) : /* @__PURE__ */ import_react27.default.createElement(OAuthProviderLogo, {
    option: props.option
  }), /* @__PURE__ */ import_react27.default.createElement("span", null, props.option.displayName));
};

// public/components/common/SignInControl.tsx
var import_react28 = __toModule(require_react());
var SignInControl = (props) => {
  const fider = useFider();
  const [email, setEmail] = (0, import_react28.useState)("");
  const [error2, setError] = (0, import_react28.useState)(void 0);
  const signIn2 = async () => {
    const result = await actions_exports.signIn(email);
    if (result.ok) {
      setEmail("");
      setError(void 0);
      if (props.onEmailSent) {
        props.onEmailSent(email);
      }
    } else if (result.error) {
      setError(result.error);
    }
  };
  const providersLen = fider.settings.oauth.length;
  if (!isCookieEnabled()) {
    return /* @__PURE__ */ import_react28.default.createElement(Message, {
      type: "error"
    }, /* @__PURE__ */ import_react28.default.createElement("h3", null, "Cookies Required"), /* @__PURE__ */ import_react28.default.createElement("p", null, "Cookies are not enabled on your browser. Please enable cookies in your browser preferences to continue."));
  }
  return /* @__PURE__ */ import_react28.default.createElement("div", {
    className: "c-signin-control"
  }, providersLen > 0 && /* @__PURE__ */ import_react28.default.createElement("div", {
    className: "l-signin-social"
  }, /* @__PURE__ */ import_react28.default.createElement("div", {
    className: "row"
  }, fider.settings.oauth.map((o, i) => /* @__PURE__ */ import_react28.default.createElement(import_react28.default.Fragment, {
    key: o.provider
  }, i % 4 === 0 && /* @__PURE__ */ import_react28.default.createElement("div", {
    className: "col-lf"
  }), /* @__PURE__ */ import_react28.default.createElement("div", {
    className: `col-sm l-provider-${o.provider} l-social-col ${providersLen === 1 ? "l-social-col-100" : ""}`
  }, /* @__PURE__ */ import_react28.default.createElement(SocialSignInButton, {
    option: o,
    redirectTo: props.redirectTo
  }))))), /* @__PURE__ */ import_react28.default.createElement("p", {
    className: "info"
  }, "We will never post to these accounts on your behalf.")), providersLen > 0 && /* @__PURE__ */ import_react28.default.createElement("div", {
    className: "c-divider"
  }, "OR"), props.useEmail && /* @__PURE__ */ import_react28.default.createElement("div", {
    className: "l-signin-email"
  }, /* @__PURE__ */ import_react28.default.createElement("p", null, "Enter your email address to sign in"), /* @__PURE__ */ import_react28.default.createElement(Form, {
    error: error2
  }, /* @__PURE__ */ import_react28.default.createElement(Input, {
    field: "email",
    value: email,
    autoFocus: !device_exports.isTouch(),
    onChange: setEmail,
    placeholder: "yourname@example.com",
    suffix: /* @__PURE__ */ import_react28.default.createElement(Button, {
      type: "submit",
      color: "positive",
      disabled: email === "",
      onClick: signIn2
    }, "Sign in")
  }))));
};

// public/components/common/Segment.tsx
var import_react29 = __toModule(require_react());
var Segment = (props) => {
  return /* @__PURE__ */ import_react29.default.createElement("div", {
    className: `c-segment ${props.className || ""}`
  }, props.children);
};

// public/components/common/List.tsx
var import_react30 = __toModule(require_react());
var List = (props) => {
  const className = classSet({
    "c-list": true,
    [props.className || ""]: true,
    "m-divided": props.divided,
    "m-hover": props.hover
  });
  return /* @__PURE__ */ import_react30.default.createElement("div", {
    className
  }, props.children);
};
var ListItem = (props) => {
  const className = classSet({
    "c-list-item": true,
    [props.className || ""]: true,
    "m-selectable": props.onClick
  });
  if (props.onClick) {
    return /* @__PURE__ */ import_react30.default.createElement("div", {
      className,
      onClick: props.onClick
    }, props.children);
  }
  return /* @__PURE__ */ import_react30.default.createElement("div", {
    className
  }, props.children);
};

// public/components/common/Moment.tsx
var import_react31 = __toModule(require_react());

// public/components/common/Modal.tsx
var import_react32 = __toModule(require_react());
var import_react_dom = __toModule(require_react_dom());
var ModalWindow = (props) => {
  const root = (0, import_react32.useRef)();
  (0, import_react32.useEffect)(() => {
    if (typeof document !== "undefined") {
      if (props.isOpen) {
        document.body.style.overflow = "hidden";
        document.addEventListener("keydown", keyDown, false);
      } else {
        document.body.style.overflow = "";
        document.removeEventListener("keydown", keyDown, false);
      }
    }
  }, [props.isOpen]);
  const swallow = (evt) => {
    evt.stopPropagation();
  };
  const keyDown = (event) => {
    if (event.keyCode === 27) {
      close();
    }
  };
  const close = () => {
    if (props.canClose) {
      props.onClose();
    }
  };
  if (!props.isOpen) {
    return null;
  }
  const className = classSet({
    "c-modal-window": true,
    [`${props.className}`]: !!props.className,
    "m-center": props.center,
    [`m-${props.size}`]: true
  });
  return import_react_dom.default.createPortal(/* @__PURE__ */ import_react32.default.createElement("div", {
    "aria-disabled": true,
    className: "c-modal-dimmer",
    onClick: close
  }, /* @__PURE__ */ import_react32.default.createElement("div", {
    className,
    onClick: swallow
  }, props.children)), root.current);
};
ModalWindow.defaultProps = {
  size: "small",
  canClose: true,
  center: true
};
var Modal = {
  Window: ModalWindow,
  Header: (props) => {
    return /* @__PURE__ */ import_react32.default.createElement("div", {
      className: "c-modal-header"
    }, props.children);
  },
  Content: (props) => {
    return /* @__PURE__ */ import_react32.default.createElement("div", {
      className: "c-modal-content"
    }, props.children);
  },
  Footer: (props) => {
    const align = props.align || "right";
    const className = classSet({
      "c-modal-footer": true,
      [`m-${align}`]: true
    });
    return /* @__PURE__ */ import_react32.default.createElement("div", {
      className
    }, props.children);
  }
};

// public/components/common/UserName.tsx
var import_react33 = __toModule(require_react());
var UserName = (props) => {
  const className = classSet({
    "c-username": true,
    "m-staff": props.user.role && isCollaborator(props.user.role)
  });
  return /* @__PURE__ */ import_react33.default.createElement("span", {
    className
  }, props.user.name || "Anonymous");
};

// public/components/common/Loader.tsx
var import_react34 = __toModule(require_react());
function Loader() {
  const [show, setShow] = (0, import_react34.useState)(false);
  useTimeout(() => {
    setShow(true);
  }, 500);
  return show ? /* @__PURE__ */ import_react34.default.createElement("div", {
    className: "c-loader"
  }) : null;
}

// public/components/common/Logo.tsx
var import_react35 = __toModule(require_react());
var TenantLogoURL = (tenant, size) => {
  if (tenant && tenant.logoBlobKey) {
    return uploadedImageURL(tenant.logoBlobKey, size);
  }
  return void 0;
};
var TenantLogo = (props) => {
  const fider = useFider();
  const tenant = fider.session.tenant;
  if (tenant && tenant.logoBlobKey) {
    return /* @__PURE__ */ import_react35.default.createElement("img", {
      src: TenantLogoURL(fider.session.tenant, props.size),
      alt: tenant.name
    });
  } else if (props.useFiderIfEmpty) {
    return /* @__PURE__ */ import_react35.default.createElement("img", {
      src: "https://getfider.com/images/logo-100x100.png",
      alt: "Fider"
    });
  }
  return null;
};
TenantLogo.defaultProps = {
  useFiderIfEmpty: false
};
var systemProvidersLogo = {
  google: `data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiA/PjwhRE9DVFlQRSBzdmcgIFBVQkxJQyAnLS8vVzNDLy9EVEQgU1ZHIDEuMS8vRU4nICAnaHR0cDovL3d3dy53My5vcmcvR3JhcGhpY3MvU1ZHLzEuMS9EVEQvc3ZnMTEuZHRkJz48c3ZnIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMCAwIDQwMCA0MDAiIGhlaWdodD0iNDAwcHgiIGlkPSJMYXllcl8xIiB2ZXJzaW9uPSIxLjEiIHZpZXdCb3g9IjAgMCA0MDAgNDAwIiB3aWR0aD0iNDAwcHgiIHhtbDpzcGFjZT0icHJlc2VydmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiPjxnPjxwYXRoIGQ9Ik0xNDIuOSwyNC4yQzk3LjYsMzkuNyw1OSw3My42LDM3LjUsMTE2LjVjLTcuNSwxNC44LTEyLjksMzAuNS0xNi4yLDQ2LjhjLTguMiw0MC40LTIuNSw4My41LDE2LjEsMTIwLjMgICBjMTIuMSwyNCwyOS41LDQ1LjQsNTAuNSw2Mi4xYzE5LjksMTUuOCw0MywyNy42LDY3LjYsMzQuMWMzMSw4LjMsNjQsOC4xLDk1LjIsMWMyOC4yLTYuNSw1NC45LTIwLDc2LjItMzkuNiAgIGMyMi41LTIwLjcsMzguNi00Ny45LDQ3LjEtNzcuMmM5LjMtMzEuOSwxMC41LTY2LDQuNy05OC44Yy01OC4zLDAtMTE2LjcsMC0xNzUsMGMwLDI0LjIsMCw0OC40LDAsNzIuNmMzMy44LDAsNjcuNiwwLDEwMS40LDAgICBjLTMuOSwyMy4yLTE3LjcsNDQuNC0zNy4yLDU3LjVjLTEyLjMsOC4zLTI2LjQsMTMuNi00MSwxNi4yYy0xNC42LDIuNS0yOS44LDIuOC00NC40LTAuMWMtMTQuOS0zLTI5LTkuMi00MS40LTE3LjkgICBjLTE5LjgtMTMuOS0zNC45LTM0LjItNDIuNi01Ny4xYy03LjktMjMuMy04LTQ5LjIsMC03Mi40YzUuNi0xNi40LDE0LjgtMzEuNSwyNy00My45YzE1LTE1LjQsMzQuNS0yNi40LDU1LjYtMzAuOSAgIGMxOC0zLjgsMzctMy4xLDU0LjYsMi4yYzE1LDQuNSwyOC44LDEyLjgsNDAuMSwyMy42YzExLjQtMTEuNCwyMi44LTIyLjgsMzQuMi0zNC4yYzYtNi4xLDEyLjMtMTIsMTguMS0xOC4zICAgYy0xNy4zLTE2LTM3LjctMjguOS01OS45LTM3LjFDMjI4LjIsMTAuNiwxODMuMiwxMC4zLDE0Mi45LDI0LjJ6IiBmaWxsPSIjRkZGRkZGIi8+PGc+PHBhdGggZD0iTTE0Mi45LDI0LjJjNDAuMi0xMy45LDg1LjMtMTMuNiwxMjUuMywxLjFjMjIuMiw4LjIsNDIuNSwyMSw1OS45LDM3LjFjLTUuOCw2LjMtMTIuMSwxMi4yLTE4LjEsMTguMyAgICBjLTExLjQsMTEuNC0yMi44LDIyLjgtMzQuMiwzNC4yYy0xMS4zLTEwLjgtMjUuMS0xOS00MC4xLTIzLjZjLTE3LjYtNS4zLTM2LjYtNi4xLTU0LjYtMi4yYy0yMSw0LjUtNDAuNSwxNS41LTU1LjYsMzAuOSAgICBjLTEyLjIsMTIuMy0yMS40LDI3LjUtMjcsNDMuOWMtMjAuMy0xNS44LTQwLjYtMzEuNS02MS00Ny4zQzU5LDczLjYsOTcuNiwzOS43LDE0Mi45LDI0LjJ6IiBmaWxsPSIjRUE0MzM1Ii8+PC9nPjxnPjxwYXRoIGQ9Ik0yMS40LDE2My4yYzMuMy0xNi4yLDguNy0zMiwxNi4yLTQ2LjhjMjAuMywxNS44LDQwLjYsMzEuNSw2MSw0Ny4zYy04LDIzLjMtOCw0OS4yLDAsNzIuNCAgICBjLTIwLjMsMTUuOC00MC42LDMxLjYtNjAuOSw0Ny4zQzE4LjksMjQ2LjcsMTMuMiwyMDMuNiwyMS40LDE2My4yeiIgZmlsbD0iI0ZCQkMwNSIvPjwvZz48Zz48cGF0aCBkPSJNMjAzLjcsMTY1LjFjNTguMywwLDExNi43LDAsMTc1LDBjNS44LDMyLjcsNC41LDY2LjgtNC43LDk4LjhjLTguNSwyOS4zLTI0LjYsNTYuNS00Ny4xLDc3LjIgICAgYy0xOS43LTE1LjMtMzkuNC0zMC42LTU5LjEtNDUuOWMxOS41LTEzLjEsMzMuMy0zNC4zLDM3LjItNTcuNWMtMzMuOCwwLTY3LjYsMC0xMDEuNCwwQzIwMy43LDIxMy41LDIwMy43LDE4OS4zLDIwMy43LDE2NS4xeiIgZmlsbD0iIzQyODVGNCIvPjwvZz48Zz48cGF0aCBkPSJNMzcuNSwyODMuNWMyMC4zLTE1LjcsNDAuNi0zMS41LDYwLjktNDcuM2M3LjgsMjIuOSwyMi44LDQzLjIsNDIuNiw1Ny4xYzEyLjQsOC43LDI2LjYsMTQuOSw0MS40LDE3LjkgICAgYzE0LjYsMywyOS43LDIuNiw0NC40LDAuMWMxNC42LTIuNiwyOC43LTcuOSw0MS0xNi4yYzE5LjcsMTUuMywzOS40LDMwLjYsNTkuMSw0NS45Yy0yMS4zLDE5LjctNDgsMzMuMS03Ni4yLDM5LjYgICAgYy0zMS4yLDcuMS02NC4yLDcuMy05NS4yLTFjLTI0LjYtNi41LTQ3LjctMTguMi02Ny42LTM0LjFDNjcsMzI4LjksNDkuNiwzMDcuNSwzNy41LDI4My41eiIgZmlsbD0iIzM0QTg1MyIvPjwvZz48L2c+PC9zdmc+`,
  facebook: `data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiA/PjwhRE9DVFlQRSBzdmcgIFBVQkxJQyAnLS8vVzNDLy9EVEQgU1ZHIDEuMC8vRU4nICAnaHR0cDovL3d3dy53My5vcmcvVFIvMjAwMS9SRUMtU1ZHLTIwMDEwOTA0L0RURC9zdmcxMC5kdGQnPjxzdmcgZW5hYmxlLWJhY2tncm91bmQ9Im5ldyAwIDAgMzIgMzIiIGhlaWdodD0iMzJweCIgaWQ9IkxheWVyXzEiIHZlcnNpb249IjEuMCIgdmlld0JveD0iMCAwIDMyIDMyIiB3aWR0aD0iMzJweCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayI+PGc+PHBhdGggZD0iTTMyLDMwYzAsMS4xMDQtMC44OTYsMi0yLDJIMmMtMS4xMDQsMC0yLTAuODk2LTItMlYyYzAtMS4xMDQsMC44OTYtMiwyLTJoMjhjMS4xMDQsMCwyLDAuODk2LDIsMlYzMHoiIGZpbGw9IiMzQjU5OTgiLz48cGF0aCBkPSJNMjIsMzJWMjBoNGwxLTVoLTV2LTJjMC0yLDEuMDAyLTMsMy0zaDJWNWMtMSwwLTIuMjQsMC00LDBjLTMuNjc1LDAtNiwyLjg4MS02LDd2M2gtNHY1aDR2MTJIMjJ6IiBmaWxsPSIjRkZGRkZGIiBpZD0iZiIvPjwvZz48Zy8+PGcvPjxnLz48Zy8+PGcvPjxnLz48L3N2Zz4=`,
  github: "data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiA/PjwhRE9DVFlQRSBzdmcgIFBVQkxJQyAnLS8vVzNDLy9EVEQgU1ZHIDEuMC8vRU4nICAnaHR0cDovL3d3dy53My5vcmcvVFIvMjAwMS9SRUMtU1ZHLTIwMDEwOTA0L0RURC9zdmcxMC5kdGQnPjxzdmcgZW5hYmxlLWJhY2tncm91bmQ9Im5ldyAwIDAgMzIgMzIiIGhlaWdodD0iMzJweCIgaWQ9IkxheWVyXzEiIHZlcnNpb249IjEuMCIgdmlld0JveD0iMCAwIDMyIDMyIiB3aWR0aD0iMzJweCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayI+PHBhdGggY2xpcC1ydWxlPSJldmVub2RkIiBkPSJNMTYuMDAzLDBDNy4xNywwLDAuMDA4LDcuMTYyLDAuMDA4LDE1Ljk5NyAgYzAsNy4wNjcsNC41ODIsMTMuMDYzLDEwLjk0LDE1LjE3OWMwLjgsMC4xNDYsMS4wNTItMC4zMjgsMS4wNTItMC43NTJjMC0wLjM4LDAuMDA4LTEuNDQyLDAtMi43NzcgIGMtNC40NDksMC45NjctNS4zNzEtMi4xMDctNS4zNzEtMi4xMDdjLTAuNzI3LTEuODQ4LTEuNzc1LTIuMzQtMS43NzUtMi4zNGMtMS40NTItMC45OTIsMC4xMDktMC45NzMsMC4xMDktMC45NzMgIGMxLjYwNSwwLjExMywyLjQ1MSwxLjY0OSwyLjQ1MSwxLjY0OWMxLjQyNywyLjQ0MywzLjc0MywxLjczNyw0LjY1NCwxLjMyOWMwLjE0Ni0xLjAzNCwwLjU2LTEuNzM5LDEuMDE3LTIuMTM5ICBjLTMuNTUyLTAuNDA0LTcuMjg2LTEuNzc2LTcuMjg2LTcuOTA2YzAtMS43NDcsMC42MjMtMy4xNzQsMS42NDYtNC4yOTJDNy4yOCwxMC40NjQsNi43Myw4LjgzNyw3LjYwMiw2LjYzNCAgYzAsMCwxLjM0My0wLjQzLDQuMzk4LDEuNjQxYzEuMjc2LTAuMzU1LDIuNjQ1LTAuNTMyLDQuMDA1LTAuNTM4YzEuMzU5LDAuMDA2LDIuNzI3LDAuMTgzLDQuMDA1LDAuNTM4ICBjMy4wNTUtMi4wNyw0LjM5Ni0xLjY0MSw0LjM5Ni0xLjY0MWMwLjg3MiwyLjIwMywwLjMyMywzLjgzLDAuMTU5LDQuMjM0YzEuMDIzLDEuMTE4LDEuNjQ0LDIuNTQ1LDEuNjQ0LDQuMjkyICBjMCw2LjE0Ni0zLjc0LDcuNDk4LTcuMzA0LDcuODkzQzE5LjQ3OSwyMy41NDgsMjAsMjQuNTA4LDIwLDI2YzAsMiwwLDMuOTAyLDAsNC40MjhjMCwwLjQyOCwwLjI1OCwwLjkwMSwxLjA3LDAuNzQ2ICBDMjcuNDIyLDI5LjA1NSwzMiwyMy4wNjIsMzIsMTUuOTk3QzMyLDcuMTYyLDI0LjgzOCwwLDE2LjAwMywweiIgZmlsbD0iIzE4MTYxNiIgZmlsbC1ydWxlPSJldmVub2RkIi8+PGcvPjxnLz48Zy8+PGcvPjxnLz48Zy8+PC9zdmc+"
};
var OAuthProviderLogoURL = (logoBlobKey) => {
  if (logoBlobKey) {
    return uploadedImageURL(logoBlobKey, 100);
  }
  return void 0;
};
var OAuthProviderLogo = (props) => {
  if (props.option.logoBlobKey) {
    return /* @__PURE__ */ import_react35.default.createElement("img", {
      src: OAuthProviderLogoURL(props.option.logoBlobKey),
      alt: props.option.displayName
    });
  }
  if (props.option.provider && props.option.provider in systemProvidersLogo) {
    return /* @__PURE__ */ import_react35.default.createElement("img", {
      src: systemProvidersLogo[props.option.provider],
      alt: props.option.displayName
    });
  }
  return null;
};

// public/components/common/Toggle.tsx
var import_react36 = __toModule(require_react());

// public/components/common/FiderVersion.tsx
var import_react37 = __toModule(require_react());

// public/components/common/DropDown.tsx
var import_react38 = __toModule(require_react());
var DropDown = class extends import_react38.default.Component {
  constructor(props) {
    super(props);
    this.mounted = false;
    this.handleMouseDown = (event) => {
      if (event.type === "mousedown" && event.button !== 0) {
        return;
      }
      event.stopPropagation();
      event.preventDefault();
      this.setState({
        isOpen: true
      }, this.addListeners);
    };
    this.renderItem = (item) => {
      if (!item) {
        return;
      }
      const {label, value} = item;
      const isSelected = this.props.highlightSelected && this.state.selected && value === this.state.selected.value;
      const className = classSet({
        "c-dropdown-item": true,
        "is-selected": isSelected
      });
      return /* @__PURE__ */ import_react38.default.createElement("div", {
        key: value,
        className,
        onMouseDown: this.setSelected.bind(this, item),
        onClick: this.setSelected.bind(this, item)
      }, item.render ? item.render : label);
    };
    this.handleDocumentClick = (event) => {
      if (this.mounted) {
        const node = this.rootElementRef.current;
        if (node && !node.contains(event.target)) {
          if (this.state.isOpen) {
            this.setState({
              isOpen: false
            }, this.removeListeners);
          }
        }
      }
    };
    this.rootElementRef = import_react38.default.createRef();
    this.state = {
      selected: this.findItem(props.defaultValue, props.items),
      isOpen: false
    };
  }
  componentDidMount() {
    this.mounted = true;
  }
  componentWillUnmount() {
    this.mounted = false;
    this.removeListeners();
  }
  addListeners() {
    document.addEventListener("click", this.handleDocumentClick, false);
    document.addEventListener("touchend", this.handleDocumentClick, false);
  }
  removeListeners() {
    document.removeEventListener("click", this.handleDocumentClick, false);
    document.removeEventListener("touchend", this.handleDocumentClick, false);
  }
  findItem(value, items) {
    for (const item of items) {
      if (item && item.value === value) {
        return item;
      }
    }
    return void 0;
  }
  setSelected(selected) {
    const newState = {
      selected,
      isOpen: false
    };
    this.fireChangeEvent(newState);
    this.setState(newState, this.removeListeners);
  }
  fireChangeEvent(newState) {
    if (newState.selected && newState.selected !== this.state.selected && this.props.onChange) {
      this.props.onChange(newState.selected);
    }
  }
  buildItemList() {
    const items = this.props.items.map(this.renderItem);
    return /* @__PURE__ */ import_react38.default.createElement("div", {
      className: "c-dropdown-menu"
    }, this.props.header && /* @__PURE__ */ import_react38.default.createElement("div", {
      className: "c-dropdown-menu-header"
    }, this.props.header), /* @__PURE__ */ import_react38.default.createElement("div", {
      className: "c-dropdown-menu-items"
    }, items.length ? items : /* @__PURE__ */ import_react38.default.createElement("div", {
      className: `c-dropdown-noresults`
    }, "No results found")));
  }
  render() {
    const text = this.state.selected ? this.state.selected.label : /* @__PURE__ */ import_react38.default.createElement("span", {
      className: "c-dropdown-placeholder"
    }, this.props.placeholder);
    const search = /* @__PURE__ */ import_react38.default.createElement("input", {
      type: "text",
      autoFocus: true,
      onChange: this.props.onSearchChange
    });
    const dropdownClass = classSet({
      "c-dropdown": true,
      [`${this.props.className}`]: this.props.className,
      "is-open": this.state.isOpen,
      [`m-style-${this.props.style}`]: true,
      "is-inline": this.props.inline,
      "m-right": this.props.direction === "right",
      "m-left": this.props.direction === "left"
    });
    return /* @__PURE__ */ import_react38.default.createElement("div", {
      ref: this.rootElementRef,
      className: dropdownClass
    }, /* @__PURE__ */ import_react38.default.createElement("div", {
      onMouseDown: this.handleMouseDown,
      onTouchEnd: this.handleMouseDown
    }, this.props.renderControl ? /* @__PURE__ */ import_react38.default.createElement("div", {
      className: "c-dropdown-control"
    }, this.props.renderControl(this.state.selected)) : /* @__PURE__ */ import_react38.default.createElement("div", {
      className: "c-dropdown-control"
    }, this.state.isOpen && this.props.searchable ? search : this.props.renderText ? this.props.renderText(this.state.selected) : /* @__PURE__ */ import_react38.default.createElement("div", null, text), /* @__PURE__ */ import_react38.default.createElement("span", {
      className: "c-dropdown-arrow"
    }))), this.state.isOpen && this.buildItemList());
  }
};
DropDown.defaultProps = {
  direction: "right",
  style: "normal",
  highlightSelected: true
};

// public/components/ShowPostResponse.tsx
var ShowPostStatus = (props) => {
  return /* @__PURE__ */ import_react39.default.createElement("span", {
    className: `status-label status-${props.status.value}`
  }, props.status.title);
};
var DuplicateDetails = (props) => {
  if (!props.response) {
    return null;
  }
  const original = props.response.original;
  if (!original) {
    return null;
  }
  return /* @__PURE__ */ import_react39.default.createElement("div", {
    className: "content"
  }, /* @__PURE__ */ import_react39.default.createElement("span", null, "\u21AA"), " ", /* @__PURE__ */ import_react39.default.createElement("a", {
    href: `/posts/${original.number}/${original.slug}`
  }, original.title));
};
var StatusDetails = (props) => {
  if (!props.response || !props.response.text) {
    return null;
  }
  return /* @__PURE__ */ import_react39.default.createElement("div", {
    className: "content"
  }, /* @__PURE__ */ import_react39.default.createElement(MultiLineText, {
    text: props.response.text,
    style: "full"
  }));
};
var ShowPostResponse = (props) => {
  const status = PostStatus.Get(props.status);
  if (props.response && (status.show || props.response.text)) {
    return /* @__PURE__ */ import_react39.default.createElement(Segment, {
      className: "l-response"
    }, status.show && /* @__PURE__ */ import_react39.default.createElement(ShowPostStatus, {
      status
    }), props.showUser && /* @__PURE__ */ import_react39.default.createElement(import_react39.default.Fragment, null, /* @__PURE__ */ import_react39.default.createElement(Avatar, {
      user: props.response.user,
      size: "small"
    }), " ", /* @__PURE__ */ import_react39.default.createElement(UserName, {
      user: props.response.user
    })), status === PostStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props));
  }
  return /* @__PURE__ */ import_react39.default.createElement("div", null);
};

// public/components/ShowTag.tsx
var import_react40 = __toModule(require_react());
var getRGB = (color) => {
  const r = color.substring(0, 2);
  const g = color.substring(2, 4);
  const b = color.substring(4, 6);
  return {
    R: parseInt(r, 16),
    G: parseInt(g, 16),
    B: parseInt(b, 16)
  };
};
var textColor = (color) => {
  const components = getRGB(color);
  const bgDelta = components.R * 0.299 + components.G * 0.587 + components.B * 0.114;
  return bgDelta > 140 ? "#333" : "#fff";
};
var ShowTag = (props) => {
  const className = classSet({
    "c-tag": true,
    [`m-${props.size || "normal"}`]: true,
    "m-circular": props.circular === true
  });
  return /* @__PURE__ */ import_react40.default.createElement("div", {
    title: `${props.tag.name}${!props.tag.isPublic ? " (Private)" : ""}`,
    className,
    style: {
      backgroundColor: `#${props.tag.color}`,
      color: textColor(props.tag.color)
    }
  }, !props.tag.isPublic && !props.circular && /* @__PURE__ */ import_react40.default.createElement(FaLock, null), props.circular ? "" : props.tag.name || "Tag");
};

// public/components/SignInModal.tsx
var import_react41 = __toModule(require_react());
var SignInModal = (props) => {
  const [confirmationAddress, setConfirmationAddress] = (0, import_react41.useState)("");
  (0, import_react41.useEffect)(() => {
    if (confirmationAddress) {
      setTimeout(() => setConfirmationAddress(""), 5e3);
    }
  }, [confirmationAddress]);
  const onEmailSent = (email) => {
    setConfirmationAddress(email);
  };
  const closeModal = () => {
    setConfirmationAddress("");
    props.onClose();
  };
  const content = confirmationAddress ? /* @__PURE__ */ import_react41.default.createElement(import_react41.default.Fragment, null, /* @__PURE__ */ import_react41.default.createElement("p", null, "We have just sent a confirmation link to ", /* @__PURE__ */ import_react41.default.createElement("b", null, confirmationAddress), ". ", /* @__PURE__ */ import_react41.default.createElement("br", null), " Click the link and you\u2019ll be signed in."), /* @__PURE__ */ import_react41.default.createElement("p", null, /* @__PURE__ */ import_react41.default.createElement("a", {
    href: "#",
    onClick: closeModal
  }, "OK"))) : /* @__PURE__ */ import_react41.default.createElement(SignInControl, {
    useEmail: true,
    onEmailSent
  });
  return /* @__PURE__ */ import_react41.default.createElement(Modal.Window, {
    isOpen: props.isOpen,
    onClose: closeModal
  }, /* @__PURE__ */ import_react41.default.createElement(Modal.Header, null, "Sign in to raise your voice"), /* @__PURE__ */ import_react41.default.createElement(Modal.Content, null, content), /* @__PURE__ */ import_react41.default.createElement(LegalFooter, null));
};

// public/components/VoteCounter.tsx
var import_react42 = __toModule(require_react());
var VoteCounter = (props) => {
  const fider = useFider();
  const [hasVoted, setHasVoted] = (0, import_react42.useState)(props.post.hasVoted);
  const [votesCount, setVotesCount] = (0, import_react42.useState)(props.post.votesCount);
  const [isSignInModalOpen, setIsSignInModalOpen] = (0, import_react42.useState)(false);
  const voteOrUndo = async () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true);
      return;
    }
    const action = hasVoted ? actions_exports.removeVote : actions_exports.addVote;
    const response = await action(props.post.number);
    if (response.ok) {
      setVotesCount(votesCount + (hasVoted ? -1 : 1));
      setHasVoted(!hasVoted);
    }
  };
  const hideModal = () => setIsSignInModalOpen(false);
  const status = PostStatus.Get(props.post.status);
  const className = classSet({
    "m-voted": !status.closed && hasVoted,
    "m-disabled": status.closed,
    "no-touch": !device_exports.isTouch()
  });
  const vote = /* @__PURE__ */ import_react42.default.createElement("button", {
    className,
    onClick: voteOrUndo
  }, /* @__PURE__ */ import_react42.default.createElement(FaCaretUp, null), votesCount);
  const disabled = /* @__PURE__ */ import_react42.default.createElement("button", {
    className
  }, /* @__PURE__ */ import_react42.default.createElement(FaCaretUp, null), votesCount);
  return /* @__PURE__ */ import_react42.default.createElement(import_react42.default.Fragment, null, /* @__PURE__ */ import_react42.default.createElement(SignInModal, {
    isOpen: isSignInModalOpen,
    onClose: hideModal
  }), /* @__PURE__ */ import_react42.default.createElement("div", {
    className: "c-vote-counter"
  }, status.closed ? disabled : vote));
};

// public/pages/Home/components/SimilarPosts.tsx
var import_react44 = __toModule(require_react());

// public/pages/Home/components/ListPosts.tsx
var import_react43 = __toModule(require_react());
var ListPostItem = (props) => {
  return /* @__PURE__ */ import_react43.default.createElement(ListItem, null, /* @__PURE__ */ import_react43.default.createElement(VoteCounter, {
    post: props.post
  }), /* @__PURE__ */ import_react43.default.createElement("div", {
    className: "c-list-item-content"
  }, props.post.commentsCount > 0 && /* @__PURE__ */ import_react43.default.createElement("div", {
    className: "info right"
  }, props.post.commentsCount, " ", /* @__PURE__ */ import_react43.default.createElement(FaRegComments, null)), /* @__PURE__ */ import_react43.default.createElement("a", {
    className: "c-list-item-title",
    href: `/posts/${props.post.number}/${props.post.slug}`
  }, props.post.title), /* @__PURE__ */ import_react43.default.createElement(MultiLineText, {
    className: "c-list-item-description",
    text: props.post.description,
    style: "simple"
  }), /* @__PURE__ */ import_react43.default.createElement(ShowPostResponse, {
    showUser: false,
    status: props.post.status,
    response: props.post.response
  }), props.tags.map((tag) => /* @__PURE__ */ import_react43.default.createElement(ShowTag, {
    key: tag.id,
    size: "tiny",
    tag
  }))));
};
var ListPosts = (props) => {
  if (!props.posts) {
    return null;
  }
  if (props.posts.length === 0) {
    return /* @__PURE__ */ import_react43.default.createElement("p", {
      className: "center"
    }, props.emptyText);
  }
  return /* @__PURE__ */ import_react43.default.createElement(List, {
    className: "c-post-list",
    divided: true
  }, props.posts.map((post) => /* @__PURE__ */ import_react43.default.createElement(ListPostItem, {
    key: post.id,
    post,
    tags: props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)
  })));
};

// public/pages/Home/components/SimilarPosts.tsx
var SimilarPosts = class extends import_react44.default.Component {
  constructor(props) {
    super(props);
    this.loadSimilarPosts = () => {
      if (this.state.loading) {
        actions_exports.searchPosts({query: this.state.title}).then((x) => {
          if (x.ok) {
            this.setState({loading: false, posts: x.data});
          }
        });
      }
    };
    this.state = {
      title: props.title,
      loading: !!props.title,
      posts: []
    };
  }
  static getDerivedStateFromProps(nextProps, prevState) {
    if (nextProps.title !== prevState.title) {
      return {
        loading: true,
        title: nextProps.title
      };
    }
    return null;
  }
  componentDidMount() {
    this.loadSimilarPosts();
  }
  componentDidUpdate() {
    window.clearTimeout(this.timer);
    this.timer = window.setTimeout(this.loadSimilarPosts, 500);
  }
  render() {
    return /* @__PURE__ */ import_react44.default.createElement(import_react44.default.Fragment, null, /* @__PURE__ */ import_react44.default.createElement(Heading, {
      title: "Similar posts",
      subtitle: "Consider voting on existing posts instead.",
      icon: FaRegLightbulb,
      size: "small",
      dividing: true
    }), this.state.loading ? /* @__PURE__ */ import_react44.default.createElement(Loader, null) : /* @__PURE__ */ import_react44.default.createElement(ListPosts, {
      posts: this.state.posts,
      tags: this.props.tags,
      emptyText: `No similar posts matched '${this.props.title}'.`
    }));
  }
};

// public/pages/Home/components/PostInput.tsx
var import_react45 = __toModule(require_react());
var CACHE_TITLE_KEY = "PostInput-Title";
var CACHE_DESCRIPTION_KEY = "PostInput-Description";
var PostInput = (props) => {
  const getCachedValue = (key) => {
    if (fider.session.isAuthenticated) {
      return cache.session.get(key) || "";
    }
    return "";
  };
  const fider = useFider();
  const titleRef = (0, import_react45.useRef)();
  const [title, setTitle] = (0, import_react45.useState)(getCachedValue(CACHE_TITLE_KEY));
  const [description, setDescription] = (0, import_react45.useState)(getCachedValue(CACHE_DESCRIPTION_KEY));
  const [isSignInModalOpen, setIsSignInModalOpen] = (0, import_react45.useState)(false);
  const [attachments, setAttachments] = (0, import_react45.useState)([]);
  const [error2, setError] = (0, import_react45.useState)(void 0);
  (0, import_react45.useEffect)(() => {
    props.onTitleChanged(title);
  }, [title]);
  const handleTitleFocus = () => {
    if (!fider.session.isAuthenticated && titleRef.current) {
      titleRef.current.blur();
      setIsSignInModalOpen(true);
    }
  };
  const handleTitleChange = (value) => {
    cache.session.set(CACHE_TITLE_KEY, value);
    setTitle(value);
    props.onTitleChanged(value);
  };
  const hideModal = () => setIsSignInModalOpen(false);
  const clearError = () => setError(void 0);
  const handleDescriptionChange = (value) => {
    cache.session.set(CACHE_DESCRIPTION_KEY, value);
    setDescription(value);
  };
  const submit = async (event) => {
    if (title) {
      const result = await actions_exports.createPost(title, description, attachments);
      if (result.ok) {
        clearError();
        cache.session.remove(CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY);
        location.href = `/posts/${result.data.number}/${result.data.slug}`;
        event.preventEnable();
      } else if (result.error) {
        setError(result.error);
      }
    }
  };
  const details = () => /* @__PURE__ */ import_react45.default.createElement(import_react45.default.Fragment, null, /* @__PURE__ */ import_react45.default.createElement(TextArea, {
    field: "description",
    onChange: handleDescriptionChange,
    value: description,
    minRows: 5,
    placeholder: "Describe your suggestion (optional)"
  }), /* @__PURE__ */ import_react45.default.createElement(MultiImageUploader, {
    field: "attachments",
    maxUploads: 3,
    previewMaxWidth: 100,
    onChange: setAttachments
  }), /* @__PURE__ */ import_react45.default.createElement(Button, {
    type: "submit",
    color: "positive",
    onClick: submit
  }, "Submit"));
  return /* @__PURE__ */ import_react45.default.createElement(import_react45.default.Fragment, null, /* @__PURE__ */ import_react45.default.createElement(SignInModal, {
    isOpen: isSignInModalOpen,
    onClose: hideModal
  }), /* @__PURE__ */ import_react45.default.createElement(Form, {
    error: error2
  }, /* @__PURE__ */ import_react45.default.createElement(Input, {
    field: "title",
    noTabFocus: !fider.session.isAuthenticated,
    inputRef: titleRef,
    onFocus: handleTitleFocus,
    maxLength: 100,
    value: title,
    onChange: handleTitleChange,
    placeholder: props.placeholder
  }), title && details()));
};

// public/pages/Home/components/PostsContainer.tsx
var import_react48 = __toModule(require_react());

// public/pages/Home/components/PostFilter.tsx
var import_react46 = __toModule(require_react());
var PostFilter = (props) => {
  const fider = useFider();
  const handleChangeView = (item) => {
    props.viewChanged(item.value);
  };
  const options = [
    {value: "trending", label: "Trending"},
    {value: "recent", label: "Recent"},
    {value: "most-wanted", label: "Most Wanted"},
    {value: "most-discussed", label: "Most Discussed"}
  ];
  if (fider.session.isAuthenticated) {
    options.push({value: "my-votes", label: "My Votes"});
  }
  PostStatus.All.filter((s) => s.filterable && props.countPerStatus[s.value]).forEach((s) => {
    options.push({
      label: s.title,
      value: s.value,
      render: /* @__PURE__ */ import_react46.default.createElement("span", null, s.title, " ", /* @__PURE__ */ import_react46.default.createElement("a", {
        className: "counter"
      }, props.countPerStatus[s.value]))
    });
  });
  const viewExists = options.filter((x) => x.value === props.activeView).length > 0;
  const activeView = viewExists ? props.activeView : "trending";
  return /* @__PURE__ */ import_react46.default.createElement("div", null, /* @__PURE__ */ import_react46.default.createElement("span", {
    className: "subtitle"
  }, "View"), /* @__PURE__ */ import_react46.default.createElement(DropDown, {
    header: "What do you want to see?",
    className: "l-post-filter",
    inline: true,
    style: "simple",
    items: options,
    defaultValue: activeView,
    onChange: handleChangeView
  }));
};

// public/pages/Home/components/TagsFilter.tsx
var import_react47 = __toModule(require_react());
var TagsFilter = class extends import_react47.default.Component {
  constructor(props) {
    super(props);
    this.onChange = (item) => {
      let selected = [];
      const idx = this.state.selected.indexOf(item.value);
      if (idx >= 0) {
        selected = this.state.selected.splice(idx, 1) && this.state.selected;
      } else {
        selected = this.state.selected.concat(item.value);
      }
      this.setState({selected});
      this.props.selectionChanged(selected);
    };
    this.renderText = () => {
      const text = this.state.selected.length === 0 ? "any tag" : this.state.selected.length === 1 ? "1 tag" : `${this.state.selected.length} tags`;
      return /* @__PURE__ */ import_react47.default.createElement(import_react47.default.Fragment, null, text);
    };
    this.state = {
      selected: props.defaultSelection
    };
  }
  render() {
    if (this.props.tags.length === 0) {
      return null;
    }
    const items = this.props.tags.map((t) => {
      return {
        value: t.slug,
        label: t.name,
        render: /* @__PURE__ */ import_react47.default.createElement("div", {
          className: this.state.selected.indexOf(t.slug) >= 0 ? "selected-tag" : ""
        }, /* @__PURE__ */ import_react47.default.createElement(FaCheck, null), /* @__PURE__ */ import_react47.default.createElement(ShowTag, {
          tag: t,
          size: "mini",
          circular: true
        }), t.name)
      };
    });
    return /* @__PURE__ */ import_react47.default.createElement("div", null, /* @__PURE__ */ import_react47.default.createElement("span", {
      className: "subtitle"
    }, "with"), /* @__PURE__ */ import_react47.default.createElement(DropDown, {
      className: "l-tags-filter",
      inline: true,
      style: "simple",
      highlightSelected: false,
      items,
      onChange: this.onChange,
      renderText: this.renderText
    }));
  }
};

// public/pages/Home/components/PostsContainer.tsx
var PostsContainer = class extends import_react48.default.Component {
  constructor(props) {
    super(props);
    this.handleViewChanged = (view) => {
      this.changeFilterCriteria({view}, true);
    };
    this.handleTagsFilterChanged = (tags) => {
      this.changeFilterCriteria({tags}, true);
    };
    this.handleSearchFilterChanged = (query) => {
      this.changeFilterCriteria({query}, true);
    };
    this.clearSearch = () => {
      this.changeFilterCriteria({query: ""}, true);
    };
    this.showMore = (event) => {
      event.preventDefault();
      this.changeFilterCriteria({limit: (this.state.limit || 30) + 10}, false);
    };
    this.getShowMoreLink = () => {
      if (this.state.posts && this.state.posts.length >= (this.state.limit || 30)) {
        return querystring_exports.set("limit", (this.state.limit || 30) + 10);
      }
    };
    this.state = {
      posts: this.props.posts,
      loading: false,
      view: querystring_exports.get("view"),
      query: querystring_exports.get("query"),
      tags: querystring_exports.getArray("tags"),
      limit: querystring_exports.getNumber("limit")
    };
  }
  changeFilterCriteria(obj, reset) {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase();
      navigator_default.replaceState(querystring_exports.stringify({
        tags: this.state.tags,
        query,
        view: this.state.view,
        limit: this.state.limit
      }));
      this.searchPosts(query, this.state.view, this.state.limit, this.state.tags, reset);
    });
  }
  async searchPosts(query, view, limit, tags, reset) {
    window.clearTimeout(this.timer);
    this.setState({posts: reset ? void 0 : this.state.posts, loading: true});
    this.timer = window.setTimeout(() => {
      actions_exports.searchPosts({query, view, limit, tags}).then((response) => {
        if (response.ok && this.state.loading) {
          this.setState({loading: false, posts: response.data});
        }
      });
    }, 500);
  }
  render() {
    const showMoreLink = this.getShowMoreLink();
    return /* @__PURE__ */ import_react48.default.createElement(import_react48.default.Fragment, null, /* @__PURE__ */ import_react48.default.createElement("div", {
      className: "row"
    }, !this.state.query && /* @__PURE__ */ import_react48.default.createElement("div", {
      className: "l-filter-col col-7 col-md-8 col-lg-9 mb-2"
    }, /* @__PURE__ */ import_react48.default.createElement(Field, null, /* @__PURE__ */ import_react48.default.createElement(PostFilter, {
      activeView: this.state.view,
      viewChanged: this.handleViewChanged,
      countPerStatus: this.props.countPerStatus
    }), /* @__PURE__ */ import_react48.default.createElement(TagsFilter, {
      tags: this.props.tags,
      selectionChanged: this.handleTagsFilterChanged,
      defaultSelection: this.state.tags
    }))), /* @__PURE__ */ import_react48.default.createElement("div", {
      className: !this.state.query ? `l-search-col col-5 col-md-4 col-lg-3 mb-2` : "col-sm-12 mb-2"
    }, /* @__PURE__ */ import_react48.default.createElement(Input, {
      field: "query",
      icon: this.state.query ? FaTimes : FaSearch,
      onIconClick: this.state.query ? this.clearSearch : void 0,
      placeholder: "Search...",
      value: this.state.query,
      onChange: this.handleSearchFilterChanged
    }))), /* @__PURE__ */ import_react48.default.createElement(ListPosts, {
      posts: this.state.posts,
      tags: this.props.tags,
      emptyText: "No results matched your search, try something different."
    }), this.state.loading && /* @__PURE__ */ import_react48.default.createElement(Loader, null), showMoreLink && /* @__PURE__ */ import_react48.default.createElement("a", {
      href: showMoreLink,
      className: "c-post-list-show-more",
      onTouchEnd: this.showMore,
      onClick: this.showMore
    }, "View more posts"));
  }
};

// public/pages/Home/Home.page.tsx
var Lonely = () => {
  const fider = useFider();
  return /* @__PURE__ */ import_react49.default.createElement("div", {
    className: "l-lonely center"
  }, /* @__PURE__ */ import_react49.default.createElement(Hint, {
    permanentCloseKey: "at-least-3-posts",
    condition: fider.session.isAuthenticated && fider.session.user.isAdministrator
  }, "It's recommended that you post ", /* @__PURE__ */ import_react49.default.createElement("strong", null, "at least 3"), " suggestions here before sharing this site. The initial content is key to start the interactions with your audience."), /* @__PURE__ */ import_react49.default.createElement("p", null, /* @__PURE__ */ import_react49.default.createElement(FaRegLightbulb, null)), /* @__PURE__ */ import_react49.default.createElement("p", null, "It's lonely out here. Start by sharing a suggestion!"));
};
var defaultWelcomeMessage = `We'd love to hear what you're thinking about. 

What can we do better? This is the place for you to vote, discuss and share ideas.`;
var HomePage = (props) => {
  const fider = useFider();
  const [title, setTitle] = (0, import_react49.useState)("");
  const isLonely = () => {
    const len = Object.keys(props.countPerStatus).length;
    if (len === 0) {
      return true;
    }
    if (len === 1 && PostStatus.Deleted.value in props.countPerStatus) {
      return true;
    }
    return false;
  };
  return /* @__PURE__ */ import_react49.default.createElement("div", {
    id: "p-home",
    className: "page container"
  }, /* @__PURE__ */ import_react49.default.createElement("div", {
    className: "row"
  }, /* @__PURE__ */ import_react49.default.createElement("div", {
    className: "l-welcome-col col-md-4"
  }, /* @__PURE__ */ import_react49.default.createElement(MultiLineText, {
    className: "welcome-message",
    text: fider.session.tenant.welcomeMessage || defaultWelcomeMessage,
    style: "full"
  }), /* @__PURE__ */ import_react49.default.createElement(PostInput, {
    placeholder: fider.session.tenant.invitation || "Enter your suggestion here...",
    onTitleChanged: setTitle
  })), /* @__PURE__ */ import_react49.default.createElement("div", {
    className: "l-posts-col col-md-8"
  }, isLonely() ? /* @__PURE__ */ import_react49.default.createElement(Lonely, null) : title ? /* @__PURE__ */ import_react49.default.createElement(SimilarPosts, {
    title,
    tags: props.tags
  }) : /* @__PURE__ */ import_react49.default.createElement(PostsContainer, {
    posts: props.posts,
    tags: props.tags,
    countPerStatus: props.countPerStatus
  }))));
};
var Home_page_default = HomePage;

// public/ssr.tsx
var import_react_icons = __toModule(require_react_icons());
function doWork(f, props) {
  let fider = Fider.initialize({...f});
  return {
    html: (0, import_server.renderToString)(/* @__PURE__ */ import_react50.default.createElement(FiderContext.Provider, {
      value: fider
    }, /* @__PURE__ */ import_react50.default.createElement(import_react_icons.IconContext.Provider, {
      value: {className: "icon"}
    }, /* @__PURE__ */ import_react50.default.createElement(Home_page_default, {
      ...props
    }))))
  };
}
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  doWork
});
