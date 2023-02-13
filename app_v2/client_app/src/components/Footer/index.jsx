import { Col, Image, Row } from "antd";
import React from "react";
import "./index.less";

import FooterGroupHeading from "./FooterGroupHeading";
import FooterGroupItem from "./FooterGroupItem";

const Footer = () => {
  return (
    <footer>
      <div className="container footer-wrapper">
        <Row gutter={16}>
          <Col span={6}>
            <div className="footer-gr">
              <FooterGroupHeading title="Chăm sóc khách hàng" />
              <FooterGroupItem name="Trung tâm trợ giúp" />
              <FooterGroupItem name="Hướng dẫn mua hàng" />
              <FooterGroupItem name="Chăm sóc khách hàng" />
              <FooterGroupItem name="Vận chuyển" />
              <FooterGroupItem name="Bảo hành" />
            </div>
          </Col>
          <Col span={6}>
            <div className="footer-gr">
              <FooterGroupHeading title="Về chúng tôi" />
              <FooterGroupItem name="Giới thiệu" />
              <FooterGroupItem name="Chính sách bảo mật" />
              <FooterGroupItem name="Điều khoản" />
              <FooterGroupItem name="Vận chuyển" />
              <FooterGroupItem name="Liên hệ" />
            </div>
          </Col>
          <Col span={6}>
            <div className="footer-gr">
              <FooterGroupHeading title="Thanh toán" />
              <Image
                className="footer-payment-logo"
                src="https://uitk14-temp-cdn-ecommerce.s3.ap-southeast-1.amazonaws.com/payment-logo/visa.png"
                preview={false}
              />
              <Image
                className="footer-payment-logo"
                src="https://uitk14-temp-cdn-ecommerce.s3.ap-southeast-1.amazonaws.com/payment-logo/mastercard.png"
                preview={false}
              />
              <Image
                className="footer-payment-logo"
                src="https://uitk14-temp-cdn-ecommerce.s3.ap-southeast-1.amazonaws.com/payment-logo/jcb.png"
                preview={false}
              />
            </div>
          </Col>
          <Col span={6}>
            <div className="footer-gr">
              <FooterGroupHeading title="Mạng xã hội" />
              <FooterGroupItem name="Facebook"></FooterGroupItem>
              <FooterGroupItem name="Instagram"></FooterGroupItem>
              <FooterGroupItem name="Youtube"></FooterGroupItem>
            </div>
          </Col>
        </Row>
      </div>
    </footer>
  );
};

export default Footer;
