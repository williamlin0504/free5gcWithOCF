export default class UEInfoDetail {
  
  ueInfoDetail = {
      amfInfo:{},
      smfInfo:{},
      ccfInfo:{}
  }

  constructor(info) {
     this.amfInfo = info.amfInfo;
     this.smfInfo = info.smfInfo;
     this.ccfInfo = info.ccfInfo;
  }
}