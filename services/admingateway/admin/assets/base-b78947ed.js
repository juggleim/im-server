import{S as d,a as _,r as u,s as v,c as f,b as s,t as i,u as t,x as b,y as S,l as w,A as x,e as a,o as A}from"./index-fd7d3cea.js";const I={class:"mb-4 app-base"},k={class:"row cim-ab-row"},y={class:"col-sm-4"},E={class:"form-control-plaintext"},N={class:"row cim-ab-row"},T={class:"col-sm-4"},g={class:"form-control-plaintext"},O={class:"row cim-ab-row"},R={class:"col-sm-4"},B={class:"cim-secret-text"},P={class:"row cim-ab-row"},V={class:"col-sm-4"},h={class:"form-control-plaintext"},C={class:"row cim-ab-row"},K={class:"col-sm-2"},D={__name:"base",setup(U){let c=w(),{currentRoute:{_rawValue:{params:{app_key:n}}}}=c;d.get(_.USER_TOKEN);let l=u({appInfo:{restricted_fields:{}},isShowSecret:!1,currentAppStatus:v.ONLINE});function r(){x.getOne({app_key:n}).then(({data:e})=>{let{cur_user_count:o,max_user_count:m}=e;a.formatProps(e,{count:"num",time:"date"}),e.use_percent=Math.floor(o/m)*100,e.expired_time=e.expired_time==-1?"永久有效":a.formatTime(e.expired_time),e.n_app_secret="********************",a.extend(l.appInfo,e)})}r();function p(){let e=!l.isShowSecret;a.extend(l,{isShowSecret:e})}return(e,o)=>(A(),f("div",I,[o[10]||(o[10]=s("ul",{class:"nav nav-underline-border ab-underline-border"},[s("li",{class:"nav-item"},[s("a",{class:"nav-link active cicon cicon-product"},"基本信息")])],-1)),s("div",k,[o[0]||(o[0]=s("label",{class:"col-sm-1 col-form-label"},"App 名称",-1)),s("div",y,[s("div",E,i(t(l).appInfo.app_name),1)]),o[1]||(o[1]=s("div",{class:"col-sm-4"},null,-1))]),s("div",N,[o[2]||(o[2]=s("label",{class:"col-sm-1 col-form-label"},"App Key",-1)),s("div",T,[s("div",g,i(t(l).appInfo.app_key),1)]),o[3]||(o[3]=s("div",{class:"col-sm-4"},null,-1))]),s("div",O,[o[4]||(o[4]=s("label",{class:"col-sm-1 col-form-label"},"App Secret",-1)),s("div",R,[s("div",{class:b(["form-control-plaintext cim-app-secret",{redfont:t(l).isShowSecret}])},[s("span",B,i(t(l).isShowSecret?t(l).appInfo.app_secret:t(l).appInfo.n_app_secret),1),s("span",{class:"cicon cicon-hide cim-secret-btn",onClick:p})],2)]),o[5]||(o[5]=s("div",{class:"col-sm-4"},null,-1))]),s("div",P,[o[6]||(o[6]=s("label",{class:"col-sm-1 col-form-label"},"App 到期时间",-1)),s("div",V,[s("div",h,i(t(l).appInfo.expired_time),1)]),o[7]||(o[7]=s("div",{class:"col-sm-4"},null,-1))]),s("div",C,[o[8]||(o[8]=s("label",{class:"col-sm-1 col-form-label"},"授权数量",-1)),s("div",K,i(t(l).appInfo.n_max_user_count==-1?"无限制":t(l).appInfo.n_max_user_count),1),o[9]||(o[9]=s("div",{class:"col-sm-4"},null,-1))]),o[11]||(o[11]=S('<div class="row cim-ab-row"><label class="col-sm-1 col-form-label">App 状态</label><div class="col-sm-4"><div class="form-control-plaintext">已上线</div></div><div class="col-sm-4"></div></div>',1))]))}};export{D as default};