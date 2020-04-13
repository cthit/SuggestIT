import styled from "styled-components";
import { DigitText, DigitDesign } from "@cthit/react-digit-components";

export const LinkText = styled(DigitText.Text)`
    &:hover {
        text-decoration: underline;
    }
    color: blue;
    cursor: pointer;
`;

export const PromptWrapper = styled(DigitDesign.Card)`
    padding: 1rem;
`;
