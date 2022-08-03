import {
  ref as t,
  watch as e,
  onMounted as n,
  openBlock as r,
  createElementBlock as i,
  createElementVNode as s,
  toDisplayString as a,
  pushScopeId as o,
  popScopeId as u,
} from "vue";
import { useApi as f } from "@directus/extensions-sdk";
"undefined" != typeof globalThis
  ? globalThis
  : "undefined" != typeof window
  ? window
  : "undefined" != typeof global
  ? global
  : "undefined" != typeof self && self;
var c = { exports: {} },
  d = (c.exports = (function () {
    var t = 1e3,
      e = 6e4,
      n = 36e5,
      r = "millisecond",
      i = "second",
      s = "minute",
      a = "hour",
      o = "day",
      u = "week",
      f = "month",
      c = "quarter",
      d = "year",
      h = "date",
      l = "Invalid Date",
      m =
        /^(\d{4})[-/]?(\d{1,2})?[-/]?(\d{0,2})[Tt\s]*(\d{1,2})?:?(\d{1,2})?:?(\d{1,2})?[.:]?(\d+)?$/,
      v =
        /\[([^\]]+)]|Y{1,4}|M{1,4}|D{1,2}|d{1,4}|H{1,2}|h{1,2}|a|A|m{1,2}|s{1,2}|Z{1,2}|SSS/g,
      $ = {
        name: "en",
        weekdays:
          "Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday".split("_"),
        months:
          "January_February_March_April_May_June_July_August_September_October_November_December".split(
            "_"
          ),
      },
      p = function (t, e, n) {
        var r = String(t);
        return !r || r.length >= e
          ? t
          : "" + Array(e + 1 - r.length).join(n) + t;
      },
      g = {
        s: p,
        z: function (t) {
          var e = -t.utcOffset(),
            n = Math.abs(e),
            r = Math.floor(n / 60),
            i = n % 60;
          return (e <= 0 ? "+" : "-") + p(r, 2, "0") + ":" + p(i, 2, "0");
        },
        m: function t(e, n) {
          if (e.date() < n.date()) return -t(n, e);
          var r = 12 * (n.year() - e.year()) + (n.month() - e.month()),
            i = e.clone().add(r, f),
            s = n - i < 0,
            a = e.clone().add(r + (s ? -1 : 1), f);
          return +(-(r + (n - i) / (s ? i - a : a - i)) || 0);
        },
        a: function (t) {
          return t < 0 ? Math.ceil(t) || 0 : Math.floor(t);
        },
        p: function (t) {
          return (
            { M: f, y: d, w: u, d: o, D: h, h: a, m: s, s: i, ms: r, Q: c }[
              t
            ] ||
            String(t || "")
              .toLowerCase()
              .replace(/s$/, "")
          );
        },
        u: function (t) {
          return void 0 === t;
        },
      },
      y = "en",
      M = {};
    M[y] = $;
    var D = function (t) {
        return t instanceof w;
      },
      S = function t(e, n, r) {
        var i;
        if (!e) return y;
        if ("string" == typeof e) {
          var s = e.toLowerCase();
          M[s] && (i = s), n && ((M[s] = n), (i = s));
          var a = e.split("-");
          if (!i && a.length > 1) return t(a[0]);
        } else {
          var o = e.name;
          (M[o] = e), (i = o);
        }
        return !r && i && (y = i), i || (!r && y);
      },
      Y = function (t, e) {
        if (D(t)) return t.clone();
        var n = "object" == typeof e ? e : {};
        return (n.date = t), (n.args = arguments), new w(n);
      },
      x = g;
    (x.l = S),
      (x.i = D),
      (x.w = function (t, e) {
        return Y(t, { locale: e.$L, utc: e.$u, x: e.$x, $offset: e.$offset });
      });
    var w = (function () {
        function $(t) {
          (this.$L = S(t.locale, null, !0)), this.parse(t);
        }
        var p = $.prototype;
        return (
          (p.parse = function (t) {
            (this.$d = (function (t) {
              var e = t.date,
                n = t.utc;
              if (null === e) return new Date(NaN);
              if (x.u(e)) return new Date();
              if (e instanceof Date) return new Date(e);
              if ("string" == typeof e && !/Z$/i.test(e)) {
                var r = e.match(m);
                if (r) {
                  var i = r[2] - 1 || 0,
                    s = (r[7] || "0").substring(0, 3);
                  return n
                    ? new Date(
                        Date.UTC(
                          r[1],
                          i,
                          r[3] || 1,
                          r[4] || 0,
                          r[5] || 0,
                          r[6] || 0,
                          s
                        )
                      )
                    : new Date(
                        r[1],
                        i,
                        r[3] || 1,
                        r[4] || 0,
                        r[5] || 0,
                        r[6] || 0,
                        s
                      );
                }
              }
              return new Date(e);
            })(t)),
              (this.$x = t.x || {}),
              this.init();
          }),
          (p.init = function () {
            var t = this.$d;
            (this.$y = t.getFullYear()),
              (this.$M = t.getMonth()),
              (this.$D = t.getDate()),
              (this.$W = t.getDay()),
              (this.$H = t.getHours()),
              (this.$m = t.getMinutes()),
              (this.$s = t.getSeconds()),
              (this.$ms = t.getMilliseconds());
          }),
          (p.$utils = function () {
            return x;
          }),
          (p.isValid = function () {
            return !(this.$d.toString() === l);
          }),
          (p.isSame = function (t, e) {
            var n = Y(t);
            return this.startOf(e) <= n && n <= this.endOf(e);
          }),
          (p.isAfter = function (t, e) {
            return Y(t) < this.startOf(e);
          }),
          (p.isBefore = function (t, e) {
            return this.endOf(e) < Y(t);
          }),
          (p.$g = function (t, e, n) {
            return x.u(t) ? this[e] : this.set(n, t);
          }),
          (p.unix = function () {
            return Math.floor(this.valueOf() / 1e3);
          }),
          (p.valueOf = function () {
            return this.$d.getTime();
          }),
          (p.startOf = function (t, e) {
            var n = this,
              r = !!x.u(e) || e,
              c = x.p(t),
              l = function (t, e) {
                var i = x.w(
                  n.$u ? Date.UTC(n.$y, e, t) : new Date(n.$y, e, t),
                  n
                );
                return r ? i : i.endOf(o);
              },
              m = function (t, e) {
                return x.w(
                  n
                    .toDate()
                    [t].apply(
                      n.toDate("s"),
                      (r ? [0, 0, 0, 0] : [23, 59, 59, 999]).slice(e)
                    ),
                  n
                );
              },
              v = this.$W,
              $ = this.$M,
              p = this.$D,
              g = "set" + (this.$u ? "UTC" : "");
            switch (c) {
              case d:
                return r ? l(1, 0) : l(31, 11);
              case f:
                return r ? l(1, $) : l(0, $ + 1);
              case u:
                var y = this.$locale().weekStart || 0,
                  M = (v < y ? v + 7 : v) - y;
                return l(r ? p - M : p + (6 - M), $);
              case o:
              case h:
                return m(g + "Hours", 0);
              case a:
                return m(g + "Minutes", 1);
              case s:
                return m(g + "Seconds", 2);
              case i:
                return m(g + "Milliseconds", 3);
              default:
                return this.clone();
            }
          }),
          (p.endOf = function (t) {
            return this.startOf(t, !1);
          }),
          (p.$set = function (t, e) {
            var n,
              u = x.p(t),
              c = "set" + (this.$u ? "UTC" : ""),
              l = ((n = {}),
              (n[o] = c + "Date"),
              (n[h] = c + "Date"),
              (n[f] = c + "Month"),
              (n[d] = c + "FullYear"),
              (n[a] = c + "Hours"),
              (n[s] = c + "Minutes"),
              (n[i] = c + "Seconds"),
              (n[r] = c + "Milliseconds"),
              n)[u],
              m = u === o ? this.$D + (e - this.$W) : e;
            if (u === f || u === d) {
              var v = this.clone().set(h, 1);
              v.$d[l](m),
                v.init(),
                (this.$d = v.set(h, Math.min(this.$D, v.daysInMonth())).$d);
            } else l && this.$d[l](m);
            return this.init(), this;
          }),
          (p.set = function (t, e) {
            return this.clone().$set(t, e);
          }),
          (p.get = function (t) {
            return this[x.p(t)]();
          }),
          (p.add = function (r, c) {
            var h,
              l = this;
            r = Number(r);
            var m = x.p(c),
              v = function (t) {
                var e = Y(l);
                return x.w(e.date(e.date() + Math.round(t * r)), l);
              };
            if (m === f) return this.set(f, this.$M + r);
            if (m === d) return this.set(d, this.$y + r);
            if (m === o) return v(1);
            if (m === u) return v(7);
            var $ = ((h = {}), (h[s] = e), (h[a] = n), (h[i] = t), h)[m] || 1,
              p = this.$d.getTime() + r * $;
            return x.w(p, this);
          }),
          (p.subtract = function (t, e) {
            return this.add(-1 * t, e);
          }),
          (p.format = function (t) {
            var e = this,
              n = this.$locale();
            if (!this.isValid()) return n.invalidDate || l;
            var r = t || "YYYY-MM-DDTHH:mm:ssZ",
              i = x.z(this),
              s = this.$H,
              a = this.$m,
              o = this.$M,
              u = n.weekdays,
              f = n.months,
              c = function (t, n, i, s) {
                return (t && (t[n] || t(e, r))) || i[n].slice(0, s);
              },
              d = function (t) {
                return x.s(s % 12 || 12, t, "0");
              },
              h =
                n.meridiem ||
                function (t, e, n) {
                  var r = t < 12 ? "AM" : "PM";
                  return n ? r.toLowerCase() : r;
                },
              m = {
                YY: String(this.$y).slice(-2),
                YYYY: this.$y,
                M: o + 1,
                MM: x.s(o + 1, 2, "0"),
                MMM: c(n.monthsShort, o, f, 3),
                MMMM: c(f, o),
                D: this.$D,
                DD: x.s(this.$D, 2, "0"),
                d: String(this.$W),
                dd: c(n.weekdaysMin, this.$W, u, 2),
                ddd: c(n.weekdaysShort, this.$W, u, 3),
                dddd: u[this.$W],
                H: String(s),
                HH: x.s(s, 2, "0"),
                h: d(1),
                hh: d(2),
                a: h(s, a, !0),
                A: h(s, a, !1),
                m: String(a),
                mm: x.s(a, 2, "0"),
                s: String(this.$s),
                ss: x.s(this.$s, 2, "0"),
                SSS: x.s(this.$ms, 3, "0"),
                Z: i,
              };
            return r.replace(v, function (t, e) {
              return e || m[t] || i.replace(":", "");
            });
          }),
          (p.utcOffset = function () {
            return 15 * -Math.round(this.$d.getTimezoneOffset() / 15);
          }),
          (p.diff = function (r, h, l) {
            var m,
              v = x.p(h),
              $ = Y(r),
              p = ($.utcOffset() - this.utcOffset()) * e,
              g = this - $,
              y = x.m(this, $);
            return (
              (y =
                ((m = {}),
                (m[d] = y / 12),
                (m[f] = y),
                (m[c] = y / 3),
                (m[u] = (g - p) / 6048e5),
                (m[o] = (g - p) / 864e5),
                (m[a] = g / n),
                (m[s] = g / e),
                (m[i] = g / t),
                m)[v] || g),
              l ? y : x.a(y)
            );
          }),
          (p.daysInMonth = function () {
            return this.endOf(f).$D;
          }),
          (p.$locale = function () {
            return M[this.$L];
          }),
          (p.locale = function (t, e) {
            if (!t) return this.$L;
            var n = this.clone(),
              r = S(t, e, !0);
            return r && (n.$L = r), n;
          }),
          (p.clone = function () {
            return x.w(this.$d, this);
          }),
          (p.toDate = function () {
            return new Date(this.valueOf());
          }),
          (p.toJSON = function () {
            return this.isValid() ? this.toISOString() : null;
          }),
          (p.toISOString = function () {
            return this.$d.toISOString();
          }),
          (p.toString = function () {
            return this.$d.toUTCString();
          }),
          $
        );
      })(),
      O = w.prototype;
    return (
      (Y.prototype = O),
      [
        ["$ms", r],
        ["$s", i],
        ["$m", s],
        ["$H", a],
        ["$W", o],
        ["$M", f],
        ["$y", d],
        ["$D", h],
      ].forEach(function (t) {
        O[t[1]] = function (e) {
          return this.$g(e, t[0], t[1]);
        };
      }),
      (Y.extend = function (t, e) {
        return t.$i || (t(e, w, Y), (t.$i = !0)), Y;
      }),
      (Y.locale = S),
      (Y.isDayjs = D),
      (Y.unix = function (t) {
        return Y(1e3 * t);
      }),
      (Y.en = M[y]),
      (Y.Ls = M),
      (Y.p = {}),
      Y
    );
  })()),
  h = { exports: {} },
  l = (h.exports = (function () {
    var t = "minute",
      e = /[+-]\d\d(?::?\d\d)?/g,
      n = /([+-]|\d\d)/g;
    return function (r, i, s) {
      var a = i.prototype;
      (s.utc = function (t) {
        return new i({ date: t, utc: !0, args: arguments });
      }),
        (a.utc = function (e) {
          var n = s(this.toDate(), { locale: this.$L, utc: !0 });
          return e ? n.add(this.utcOffset(), t) : n;
        }),
        (a.local = function () {
          return s(this.toDate(), { locale: this.$L, utc: !1 });
        });
      var o = a.parse;
      a.parse = function (t) {
        t.utc && (this.$u = !0),
          this.$utils().u(t.$offset) || (this.$offset = t.$offset),
          o.call(this, t);
      };
      var u = a.init;
      a.init = function () {
        if (this.$u) {
          var t = this.$d;
          (this.$y = t.getUTCFullYear()),
            (this.$M = t.getUTCMonth()),
            (this.$D = t.getUTCDate()),
            (this.$W = t.getUTCDay()),
            (this.$H = t.getUTCHours()),
            (this.$m = t.getUTCMinutes()),
            (this.$s = t.getUTCSeconds()),
            (this.$ms = t.getUTCMilliseconds());
        } else u.call(this);
      };
      var f = a.utcOffset;
      a.utcOffset = function (r, i) {
        var s = this.$utils().u;
        if (s(r))
          return this.$u ? 0 : s(this.$offset) ? f.call(this) : this.$offset;
        if (
          "string" == typeof r &&
          ((r = (function (t) {
            void 0 === t && (t = "");
            var r = t.match(e);
            if (!r) return null;
            var i = ("" + r[0]).match(n) || ["-", 0, 0],
              s = i[0],
              a = 60 * +i[1] + +i[2];
            return 0 === a ? 0 : "+" === s ? a : -a;
          })(r)),
          null === r)
        )
          return this;
        var a = Math.abs(r) <= 16 ? 60 * r : r,
          o = this;
        if (i) return (o.$offset = a), (o.$u = 0 === r), o;
        if (0 !== r) {
          var u = this.$u
            ? this.toDate().getTimezoneOffset()
            : -1 * this.utcOffset();
          ((o = this.local().add(a + u, t)).$offset = a),
            (o.$x.$localOffset = u);
        } else o = this.utc();
        return o;
      };
      var c = a.format;
      (a.format = function (t) {
        var e = t || (this.$u ? "YYYY-MM-DDTHH:mm:ss[Z]" : "");
        return c.call(this, e);
      }),
        (a.valueOf = function () {
          var t = this.$utils().u(this.$offset)
            ? 0
            : this.$offset +
              (this.$x.$localOffset || new Date().getTimezoneOffset());
          return this.$d.valueOf() - 6e4 * t;
        }),
        (a.isUTC = function () {
          return !!this.$u;
        }),
        (a.toISOString = function () {
          return this.toDate().toISOString();
        }),
        (a.toString = function () {
          return this.toDate().toUTCString();
        });
      var d = a.toDate;
      a.toDate = function (t) {
        return "s" === t && this.$offset
          ? s(this.format("YYYY-MM-DD HH:mm:ss:SSS")).toDate()
          : d.call(this);
      };
      var h = a.diff;
      a.diff = function (t, e, n) {
        if (t && this.$u === t.$u) return h.call(this, t, e, n);
        var r = this.local(),
          i = s(t).local();
        return h.call(r, i, e, n);
      };
    };
  })()),
  m = { exports: {} },
  v = (m.exports = (function () {
    var t = { year: 0, month: 1, day: 2, hour: 3, minute: 4, second: 5 },
      e = {};
    return function (n, r, i) {
      var s,
        a = function (t, n, r) {
          void 0 === r && (r = {});
          var i = new Date(t),
            s = (function (t, n) {
              void 0 === n && (n = {});
              var r = n.timeZoneName || "short",
                i = t + "|" + r,
                s = e[i];
              return (
                s ||
                  ((s = new Intl.DateTimeFormat("en-US", {
                    hour12: !1,
                    timeZone: t,
                    year: "numeric",
                    month: "2-digit",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                    timeZoneName: r,
                  })),
                  (e[i] = s)),
                s
              );
            })(n, r);
          return s.formatToParts(i);
        },
        o = function (e, n) {
          for (var r = a(e, n), s = [], o = 0; o < r.length; o += 1) {
            var u = r[o],
              f = u.type,
              c = u.value,
              d = t[f];
            d >= 0 && (s[d] = parseInt(c, 10));
          }
          var h = s[3],
            l = 24 === h ? 0 : h,
            m =
              s[0] +
              "-" +
              s[1] +
              "-" +
              s[2] +
              " " +
              l +
              ":" +
              s[4] +
              ":" +
              s[5] +
              ":000",
            v = +e;
          return (i.utc(m).valueOf() - (v -= v % 1e3)) / 6e4;
        },
        u = r.prototype;
      (u.tz = function (t, e) {
        void 0 === t && (t = s);
        var n = this.utcOffset(),
          r = this.toDate(),
          a = r.toLocaleString("en-US", { timeZone: t }),
          o = Math.round((r - new Date(a)) / 1e3 / 60),
          u = i(a)
            .$set("millisecond", this.$ms)
            .utcOffset(15 * -Math.round(r.getTimezoneOffset() / 15) - o, !0);
        if (e) {
          var f = u.utcOffset();
          u = u.add(n - f, "minute");
        }
        return (u.$x.$timezone = t), u;
      }),
        (u.offsetName = function (t) {
          var e = this.$x.$timezone || i.tz.guess(),
            n = a(this.valueOf(), e, { timeZoneName: t }).find(function (t) {
              return "timezonename" === t.type.toLowerCase();
            });
          return n && n.value;
        });
      var f = u.startOf;
      (u.startOf = function (t, e) {
        if (!this.$x || !this.$x.$timezone) return f.call(this, t, e);
        var n = i(this.format("YYYY-MM-DD HH:mm:ss:SSS"));
        return f.call(n, t, e).tz(this.$x.$timezone, !0);
      }),
        (i.tz = function (t, e, n) {
          var r = n && e,
            a = n || e || s,
            u = o(+i(), a);
          if ("string" != typeof t) return i(t).tz(a);
          var f = (function (t, e, n) {
              var r = t - 60 * e * 1e3,
                i = o(r, n);
              if (e === i) return [r, e];
              var s = o((r -= 60 * (i - e) * 1e3), n);
              return i === s
                ? [r, i]
                : [t - 60 * Math.min(i, s) * 1e3, Math.max(i, s)];
            })(i.utc(t, r).valueOf(), u, a),
            c = f[0],
            d = f[1],
            h = i(c).utcOffset(d);
          return (h.$x.$timezone = a), h;
        }),
        (i.tz.guess = function () {
          return Intl.DateTimeFormat().resolvedOptions().timeZone;
        }),
        (i.tz.setDefault = function (t) {
          s = t;
        });
    };
  })()),
  $ = { exports: {} },
  p = ($.exports = (function () {
    var t = {
        LTS: "h:mm:ss A",
        LT: "h:mm A",
        L: "MM/DD/YYYY",
        LL: "MMMM D, YYYY",
        LLL: "MMMM D, YYYY h:mm A",
        LLLL: "dddd, MMMM D, YYYY h:mm A",
      },
      e =
        /(\[[^[]*\])|([-:/.()\s]+)|(A|a|YYYY|YY?|MM?M?M?|Do|DD?|hh?|HH?|mm?|ss?|S{1,3}|z|ZZ?)/g,
      n = /\d\d/,
      r = /\d\d?/,
      i = /\d*[^\s\d-_:/()]+/,
      s = {},
      a = function (t) {
        return (t = +t) + (t > 68 ? 1900 : 2e3);
      },
      o = function (t) {
        return function (e) {
          this[t] = +e;
        };
      },
      u = [
        /[+-]\d\d:?(\d\d)?|Z/,
        function (t) {
          (this.zone || (this.zone = {})).offset = (function (t) {
            if (!t) return 0;
            if ("Z" === t) return 0;
            var e = t.match(/([+-]|\d\d)/g),
              n = 60 * e[1] + (+e[2] || 0);
            return 0 === n ? 0 : "+" === e[0] ? -n : n;
          })(t);
        },
      ],
      f = function (t) {
        var e = s[t];
        return e && (e.indexOf ? e : e.s.concat(e.f));
      },
      c = function (t, e) {
        var n,
          r = s.meridiem;
        if (r) {
          for (var i = 1; i <= 24; i += 1)
            if (t.indexOf(r(i, 0, e)) > -1) {
              n = i > 12;
              break;
            }
        } else n = t === (e ? "pm" : "PM");
        return n;
      },
      d = {
        A: [
          i,
          function (t) {
            this.afternoon = c(t, !1);
          },
        ],
        a: [
          i,
          function (t) {
            this.afternoon = c(t, !0);
          },
        ],
        S: [
          /\d/,
          function (t) {
            this.milliseconds = 100 * +t;
          },
        ],
        SS: [
          n,
          function (t) {
            this.milliseconds = 10 * +t;
          },
        ],
        SSS: [
          /\d{3}/,
          function (t) {
            this.milliseconds = +t;
          },
        ],
        s: [r, o("seconds")],
        ss: [r, o("seconds")],
        m: [r, o("minutes")],
        mm: [r, o("minutes")],
        H: [r, o("hours")],
        h: [r, o("hours")],
        HH: [r, o("hours")],
        hh: [r, o("hours")],
        D: [r, o("day")],
        DD: [n, o("day")],
        Do: [
          i,
          function (t) {
            var e = s.ordinal,
              n = t.match(/\d+/);
            if (((this.day = n[0]), e))
              for (var r = 1; r <= 31; r += 1)
                e(r).replace(/\[|\]/g, "") === t && (this.day = r);
          },
        ],
        M: [r, o("month")],
        MM: [n, o("month")],
        MMM: [
          i,
          function (t) {
            var e = f("months"),
              n =
                (
                  f("monthsShort") ||
                  e.map(function (t) {
                    return t.slice(0, 3);
                  })
                ).indexOf(t) + 1;
            if (n < 1) throw new Error();
            this.month = n % 12 || n;
          },
        ],
        MMMM: [
          i,
          function (t) {
            var e = f("months").indexOf(t) + 1;
            if (e < 1) throw new Error();
            this.month = e % 12 || e;
          },
        ],
        Y: [/[+-]?\d+/, o("year")],
        YY: [
          n,
          function (t) {
            this.year = a(t);
          },
        ],
        YYYY: [/\d{4}/, o("year")],
        Z: u,
        ZZ: u,
      };
    function h(n) {
      var r, i;
      (r = n), (i = s && s.formats);
      for (
        var a = (n = r.replace(
            /(\[[^\]]+])|(LTS?|l{1,4}|L{1,4})/g,
            function (e, n, r) {
              var s = r && r.toUpperCase();
              return (
                n ||
                i[r] ||
                t[r] ||
                i[s].replace(
                  /(\[[^\]]+])|(MMMM|MM|DD|dddd)/g,
                  function (t, e, n) {
                    return e || n.slice(1);
                  }
                )
              );
            }
          )).match(e),
          o = a.length,
          u = 0;
        u < o;
        u += 1
      ) {
        var f = a[u],
          c = d[f],
          h = c && c[0],
          l = c && c[1];
        a[u] = l ? { regex: h, parser: l } : f.replace(/^\[|\]$/g, "");
      }
      return function (t) {
        for (var e = {}, n = 0, r = 0; n < o; n += 1) {
          var i = a[n];
          if ("string" == typeof i) r += i.length;
          else {
            var s = i.regex,
              u = i.parser,
              f = t.slice(r),
              c = s.exec(f)[0];
            u.call(e, c), (t = t.replace(c, ""));
          }
        }
        return (
          (function (t) {
            var e = t.afternoon;
            if (void 0 !== e) {
              var n = t.hours;
              e ? n < 12 && (t.hours += 12) : 12 === n && (t.hours = 0),
                delete t.afternoon;
            }
          })(e),
          e
        );
      };
    }
    return function (t, e, n) {
      (n.p.customParseFormat = !0),
        t && t.parseTwoDigitYear && (a = t.parseTwoDigitYear);
      var r = e.prototype,
        i = r.parse;
      r.parse = function (t) {
        var e = t.date,
          r = t.utc,
          a = t.args;
        this.$u = r;
        var o = a[1];
        if ("string" == typeof o) {
          var u = !0 === a[2],
            f = !0 === a[3],
            c = u || f,
            d = a[2];
          f && (d = a[2]),
            (s = this.$locale()),
            !u && d && (s = n.Ls[d]),
            (this.$d = (function (t, e, n) {
              try {
                if (["x", "X"].indexOf(e) > -1)
                  return new Date(("X" === e ? 1e3 : 1) * t);
                var r = h(e)(t),
                  i = r.year,
                  s = r.month,
                  a = r.day,
                  o = r.hours,
                  u = r.minutes,
                  f = r.seconds,
                  c = r.milliseconds,
                  d = r.zone,
                  l = new Date(),
                  m = a || (i || s ? 1 : l.getDate()),
                  v = i || l.getFullYear(),
                  $ = 0;
                (i && !s) || ($ = s > 0 ? s - 1 : l.getMonth());
                var p = o || 0,
                  g = u || 0,
                  y = f || 0,
                  M = c || 0;
                return d
                  ? new Date(
                      Date.UTC(v, $, m, p, g, y, M + 60 * d.offset * 1e3)
                    )
                  : n
                  ? new Date(Date.UTC(v, $, m, p, g, y, M))
                  : new Date(v, $, m, p, g, y, M);
              } catch (t) {
                return new Date("");
              }
            })(e, o, r)),
            this.init(),
            d && !0 !== d && (this.$L = this.locale(d).$L),
            c && e != this.format(o) && (this.$d = new Date("")),
            (s = {});
        } else if (o instanceof Array)
          for (var l = o.length, m = 1; m <= l; m += 1) {
            a[1] = o[m - 1];
            var v = n.apply(this, a);
            if (v.isValid()) {
              (this.$d = v.$d), (this.$L = v.$L), this.init();
              break;
            }
            m === l && (this.$d = new Date(""));
          }
        else i.call(this, t);
      };
    };
  })());
d.extend(l), d.extend(v), d.extend(p), d.tz.setDefault("America/Los_Angeles");
var g = {
  props: {
    day: { type: Date, default: "2021-08-31" },
    fcfsId: { type: String, default: "", required: !0 },
    tableName: { type: String, default: "durations", required: !0 },
  },
  setup(r) {
    var i = t(0),
      s = t(0),
      a = t(0),
      o = t(!0);
    const u = f();
    return (
      e([() => r.day], () => {
        c();
      }),
      n(c),
      {
        totalDevices: i,
        totalPatrons: s,
        totalMinutes: a,
        isLoading: o,
        sensorID: r.fcfsId,
        previousDay: function (t) {
          return d(t).subtract(1, "day").format("YYYY-MM-DD");
        },
        previousWeek: function (t) {
          return d(t).subtract(1, "week").format("YYYY-MM-DD");
        },
        nextWeek: l,
      }
    );
    async function c() {
      o.value = !0;
      const t = {
          _and: [
            { start: { _gte: h(r.day) } },
            { end: { _lt: h(l(r.day)) } },
            { fcfs_seq_id: { _eq: r.fcfsId } },
          ],
        },
        e = await u.get(`/items/${r.tableName}`, {
          params: { aggregate: { count: "*" }, filter: t },
        });
      i.value = parseInt(e.data.data[0].count).toLocaleString();
      const n = await u.get(`/items/${r.tableName}`, {
        params: { aggregate: { countDistinct: "patron_index" }, filter: t },
      });
      s.value = parseInt(
        n.data.data[0].countDistinct.patron_index
      ).toLocaleString();
      const f = await u.get(`/items/${r.tableName}`, {
          params: { aggregate: { sum: ["end", "start"] }, filter: t },
        }),
        c = f.data.data[0].sum.end - f.data.data[0].sum.start;
      (a.value = Math.trunc(c / 60).toLocaleString()), (o.value = !1);
    }
    function h(t) {
      return d(t).unix();
    }
    function l(t) {
      return d(t).add(1, "week").format("YYYY-MM-DD");
    }
  },
};
const y = (t) => (o("data-v-0129a3fb"), (t = t()), u(), t),
  M = { class: "header" },
  D = { class: "header" },
  S = { key: 0, class: "body" },
  Y = [y(() => s("p", { class: "text" }, " Loading... ", -1))],
  x = { key: 1, class: "body" },
  w = { class: "text" },
  O = { class: "text" },
  b = { class: "text" },
  T = { class: "footer" };
var L = [],
  _ = [];
!(function (t, e) {
  if (t && "undefined" != typeof document) {
    var n,
      r = !0 === e.prepend ? "prepend" : "append",
      i = !0 === e.singleTag,
      s =
        "string" == typeof e.container
          ? document.querySelector(e.container)
          : document.getElementsByTagName("head")[0];
    if (i) {
      var a = L.indexOf(s);
      -1 === a && ((a = L.push(s) - 1), (_[a] = {})),
        (n = _[a] && _[a][r] ? _[a][r] : (_[a][r] = o()));
    } else n = o();
    65279 === t.charCodeAt(0) && (t = t.substring(1)),
      n.styleSheet
        ? (n.styleSheet.cssText += t)
        : n.appendChild(document.createTextNode(t));
  }
  function o() {
    var t = document.createElement("style");
    if ((t.setAttribute("type", "text/css"), e.attributes))
      for (var n = Object.keys(e.attributes), i = 0; i < n.length; i++)
        t.setAttribute(n[i], e.attributes[n[i]]);
    var a = "prepend" === r ? "afterbegin" : "beforeend";
    return s.insertAdjacentElement(a, t), t;
  }
})(
  '\n.header[data-v-0129a3fb] {\n\t display: flex;\n\t margin: 0 1rem;\n\t font-weight: bold;\n}\n.text[data-v-0129a3fb] {\n\t padding: 8px 0;\n}\n.body[data-v-0129a3fb] {\n\t margin: 1rem;\n\t min-height: 9rem; /* prevent "load" flicker */\n}\n.footer[data-v-0129a3fb] {\n\t margin: 0 1rem;\n\t display: flex;\n\t flex-direction: row;\n\t justify-content: space-between;\n}\n.button[data-v-0129a3fb] {\n\t border: 1px solid #777;\n\t border-radius: 5px;\n\t cursor: pointer;\n\t padding: 0.25rem 0.5rem;\n\t background: #efefef;\n}\n',
  {}
),
  (g.render = function (t, e, n, o, u, f) {
    return (
      r(),
      i("div", null, [
        s("h1", M, " Weekly sessions for sensor " + a(n.fcfsId), 1),
        s(
          "h2",
          D,
          a(n.day) + " through " + a(o.previousDay(o.nextWeek(n.day))),
          1
        ),
        o.isLoading
          ? (r(), i("div", S, Y))
          : (r(),
            i("div", x, [
              s("p", w, a(o.totalDevices) + " devices seen ", 1),
              s("p", O, a(o.totalPatrons) + " patrons ", 1),
              s("p", b, a(o.totalMinutes) + " minutes served ", 1),
            ])),
        s("div", T, [
          s(
            "div",
            {
              class: "button",
              onClick: e[0] || (e[0] = (t) => (n.day = o.previousWeek(n.day))),
            },
            " ← Previous week "
          ),
          s(
            "div",
            {
              class: "button",
              onClick: e[1] || (e[1] = (t) => (n.day = o.nextWeek(n.day))),
            },
            " Next week → "
          ),
        ]),
      ])
    );
  }),
  (g.__scopeId = "data-v-0129a3fb"),
  (g.__file = "src/panel.vue");
var k = {
  id: "single-sensor-sessions-single-week",
  name: "Sessions/Week per sensor",
  icon: "calendar_month",
  description: "Single Sensor Sessions in a Single Week",
  component: g,
  options: ({ options: t }) => (
    console.log(t),
    [
      {
        field: "tableName",
        type: "string",
        name: "Collection",
        meta: {
          interface: "system-collection",
          options: { includeSystem: !1 },
        },
      },
      {
        field: "fcfsId",
        name: "Sensor name",
        type: "string",
        meta: {
          interface: "select-dropdown",
          options: {
            choices: [
              { text: "springfield", value: "ME8675-309" },
              { text: "in-op", value: "GA0027-004" },
              { text: "rpi03", value: "GA0058-005" },
            ],
          },
        },
      },
    ]
  ),
  minWidth: 20,
  minHeight: 15,
};
export { k as default };
