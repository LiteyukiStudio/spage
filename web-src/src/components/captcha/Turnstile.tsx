import { CaptchaProps } from "@/types/captcha";
import { Turnstile } from "@marsidev/react-turnstile";

export default function TurnstileWidget(props: CaptchaProps) {
  return <Turnstile siteKey={props.siteKey} onSuccess={props.onSuccess}/>;
}
