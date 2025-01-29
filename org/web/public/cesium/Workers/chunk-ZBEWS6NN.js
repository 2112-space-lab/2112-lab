/**
 * @license
 * Cesium - https://github.com/CesiumGS/cesium
 * Version 1.125
 *
 * Copyright 2011-2022 Cesium Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Columbus View (Pat. Pend.)
 *
 * Portions licensed separately.
 * See https://github.com/CesiumGS/cesium/blob/main/LICENSE.md for full licensing details.
 */

import{a as p,b as P,d as q}from"./chunk-RH3GFHG2.js";import{a}from"./chunk-FRWNWNYJ.js";import{a as T}from"./chunk-TA3RE4KQ.js";import{a as I,b as g}from"./chunk-RTY3VPG6.js";import{e as l}from"./chunk-LRNH5AEO.js";function y(n,i,o){if(n===0)return i*o;let t=n*n,e=t*t,c=e*t,r=c*t,h=r*t,s=h*t,d=o,u=Math.sin(2*d),f=Math.sin(4*d),M=Math.sin(6*d),_=Math.sin(8*d),E=Math.sin(10*d),S=Math.sin(12*d);return i*((1-t/4-3*e/64-5*c/256-175*r/16384-441*h/65536-4851*s/1048576)*d-(3*t/8+3*e/32+45*c/1024+105*r/4096+2205*h/131072+6237*s/524288)*u+(15*e/256+45*c/1024+525*r/16384+1575*h/65536+155925*s/8388608)*f-(35*c/3072+175*r/12288+3675*h/262144+13475*s/1048576)*M+(315*r/131072+2205*h/524288+43659*s/8388608)*_-(693*h/1310720+6237*s/5242880)*E+1001*s/8388608*S)}function z(n,i,o){let t=n/o;if(i===0)return t;let e=t*t,c=e*t,r=c*t,h=i,s=h*h,d=s*s,u=d*s,f=u*s,M=f*s,_=M*s,E=Math.sin(2*t),S=Math.cos(2*t),W=Math.sin(4*t),V=Math.cos(4*t),C=Math.sin(6*t),N=Math.cos(6*t),R=Math.sin(8*t),b=Math.cos(8*t),x=Math.sin(10*t),U=Math.cos(10*t),H=Math.sin(12*t);return t+t*s/4+7*t*d/64+15*t*u/256+579*t*f/16384+1515*t*M/65536+16837*t*_/1048576+(3*t*d/16+45*t*u/256-t*(32*e-561)*f/4096-t*(232*e-1677)*M/16384+t*(399985-90560*e+512*r)*_/5242880)*S+(21*t*u/256+483*t*f/4096-t*(224*e-1969)*M/16384-t*(33152*e-112599)*_/1048576)*V+(151*t*f/4096+4681*t*M/65536+1479*t*_/16384-453*c*_/32768)*N+(1097*t*M/65536+42783*t*_/1048576)*b+8011*t*_/1048576*U+(3*s/8+3*d/16+213*u/2048-3*e*u/64+255*f/4096-33*e*f/512+20861*M/524288-33*e*M/512+r*M/1024+28273*_/1048576-471*e*_/8192+9*r*_/4096)*E+(21*d/256+21*u/256+533*f/8192-21*e*f/512+197*M/4096-315*e*M/4096+584039*_/16777216-12517*e*_/131072+7*r*_/2048)*W+(151*u/6144+151*f/4096+5019*M/131072-453*e*M/16384+26965*_/786432-8607*e*_/131072)*C+(1097*f/131072+1097*M/65536+225797*_/10485760-1097*e*_/65536)*R+(8011*M/2621440+8011*_/1048576)*x+293393*_/251658240*H}function O(n,i){if(n===0)return Math.log(Math.tan(.5*(a.PI_OVER_TWO+i)));let o=n*Math.sin(i);return Math.log(Math.tan(.5*(a.PI_OVER_TWO+i)))-n/2*Math.log((1+o)/(1-o))}function k(n,i,o,t,e){let c=O(n._ellipticity,o),r=O(n._ellipticity,e);return Math.atan2(a.negativePiToPi(t-i),r-c)}function A(n,i,o,t,e,c,r){let h=n._heading,s=c-t,d=0;if(a.equalsEpsilon(Math.abs(h),a.PI_OVER_TWO,a.EPSILON8))if(i===o)d=i*Math.cos(e)*a.negativePiToPi(s);else{let u=Math.sin(e);d=i*Math.cos(e)*a.negativePiToPi(s)/Math.sqrt(1-n._ellipticitySquared*u*u)}else{let u=y(n._ellipticity,i,e);d=(y(n._ellipticity,i,r)-u)/Math.cos(h)}return Math.abs(d)}var B=new p,w=new p;function D(n,i,o,t){let e=p.normalize(t.cartographicToCartesian(i,w),B),c=p.normalize(t.cartographicToCartesian(o,w),w);g.typeOf.number.greaterThanOrEquals("value",Math.abs(Math.abs(p.angleBetween(e,c))-Math.PI),.0125);let r=t.maximumRadius,h=t.minimumRadius,s=r*r,d=h*h;n._ellipticitySquared=(s-d)/s,n._ellipticity=Math.sqrt(n._ellipticitySquared),n._start=P.clone(i,n._start),n._start.height=0,n._end=P.clone(o,n._end),n._end.height=0,n._heading=k(n,i.longitude,i.latitude,o.longitude,o.latitude),n._distance=A(n,t.maximumRadius,t.minimumRadius,i.longitude,i.latitude,o.longitude,o.latitude)}function v(n,i,o,t,e,c){if(o===0)return P.clone(n,c);let r=e*e,h,s,d;if(Math.abs(a.PI_OVER_TWO-Math.abs(i))>a.EPSILON8){let u=y(e,t,n.latitude),f=o*Math.cos(i),M=u+f;if(s=z(M,e,t),Math.abs(i)<a.EPSILON10)h=a.negativePiToPi(n.longitude);else{let _=O(e,n.latitude),E=O(e,s);d=Math.tan(i)*(E-_),h=a.negativePiToPi(n.longitude+d)}}else{s=n.latitude;let u;if(e===0)u=t*Math.cos(n.latitude);else{let f=Math.sin(n.latitude);u=t*Math.cos(n.latitude)/Math.sqrt(1-r*f*f)}d=o/u,i>0?h=a.negativePiToPi(n.longitude+d):h=a.negativePiToPi(n.longitude-d)}return l(c)?(c.longitude=h,c.latitude=s,c.height=0,c):new P(h,s,0)}function m(n,i,o){let t=T(o,q.default);this._ellipsoid=t,this._start=new P,this._end=new P,this._heading=void 0,this._distance=void 0,this._ellipticity=void 0,this._ellipticitySquared=void 0,l(n)&&l(i)&&D(this,n,i,t)}Object.defineProperties(m.prototype,{ellipsoid:{get:function(){return this._ellipsoid}},surfaceDistance:{get:function(){return g.defined("distance",this._distance),this._distance}},start:{get:function(){return this._start}},end:{get:function(){return this._end}},heading:{get:function(){return g.defined("distance",this._distance),this._heading}}});m.fromStartHeadingDistance=function(n,i,o,t,e){g.defined("start",n),g.defined("heading",i),g.defined("distance",o),g.typeOf.number.greaterThan("distance",o,0);let c=T(t,q.default),r=c.maximumRadius,h=c.minimumRadius,s=r*r,d=h*h,u=Math.sqrt((s-d)/s);i=a.negativePiToPi(i);let f=v(n,i,o,c.maximumRadius,u);return!l(e)||l(t)&&!t.equals(e.ellipsoid)?new m(n,f,c):(e.setEndPoints(n,f),e)};m.prototype.setEndPoints=function(n,i){g.defined("start",n),g.defined("end",i),D(this,n,i,this._ellipsoid)};m.prototype.interpolateUsingFraction=function(n,i){return this.interpolateUsingSurfaceDistance(n*this._distance,i)};m.prototype.interpolateUsingSurfaceDistance=function(n,i){if(g.typeOf.number("distance",n),!l(this._distance)||this._distance===0)throw new I("EllipsoidRhumbLine must have distinct start and end set.");return v(this._start,this._heading,n,this._ellipsoid.maximumRadius,this._ellipticity,i)};m.prototype.findIntersectionWithLongitude=function(n,i){if(g.typeOf.number("intersectionLongitude",n),!l(this._distance)||this._distance===0)throw new I("EllipsoidRhumbLine must have distinct start and end set.");let o=this._ellipticity,t=this._heading,e=Math.abs(t),c=this._start;if(n=a.negativePiToPi(n),a.equalsEpsilon(Math.abs(n),Math.PI,a.EPSILON14)&&(n=a.sign(c.longitude)*Math.PI),l(i)||(i=new P),Math.abs(a.PI_OVER_TWO-e)<=a.EPSILON8)return i.longitude=n,i.latitude=c.latitude,i.height=0,i;if(a.equalsEpsilon(Math.abs(a.PI_OVER_TWO-e),a.PI_OVER_TWO,a.EPSILON8))return a.equalsEpsilon(n,c.longitude,a.EPSILON12)?void 0:(i.longitude=n,i.latitude=a.PI_OVER_TWO*a.sign(a.PI_OVER_TWO-t),i.height=0,i);let r=c.latitude,h=o*Math.sin(r),s=Math.tan(.5*(a.PI_OVER_TWO+r))*Math.exp((n-c.longitude)/Math.tan(t)),d=(1+h)/(1-h),u=c.latitude,f;do{f=u;let M=o*Math.sin(f),_=(1+M)/(1-M);u=2*Math.atan(s*Math.pow(_/d,o/2))-a.PI_OVER_TWO}while(!a.equalsEpsilon(u,f,a.EPSILON12));return i.longitude=n,i.latitude=u,i.height=0,i};m.prototype.findIntersectionWithLatitude=function(n,i){if(g.typeOf.number("intersectionLatitude",n),!l(this._distance)||this._distance===0)throw new I("EllipsoidRhumbLine must have distinct start and end set.");let o=this._ellipticity,t=this._heading,e=this._start;if(a.equalsEpsilon(Math.abs(t),a.PI_OVER_TWO,a.EPSILON8))return;let c=O(o,e.latitude),r=O(o,n),h=Math.tan(t)*(r-c),s=a.negativePiToPi(e.longitude+h);return l(i)?(i.longitude=s,i.latitude=n,i.height=0,i):new P(s,n,0)};var $=m;export{$ as a};
