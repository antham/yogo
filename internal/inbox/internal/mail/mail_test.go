package mail

import (
	"testing"

	"github.com/antham/yogo/v3/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	mail, err := Parse[client.MailHTMLDoc](getDoc[client.MailHTMLDoc](t, "html_mail.html"))
	assert.NoError(t, err)

	content, err := mail.Coloured()
	assert.NoError(t, err)
	assert.Equal(t, `---
From  : Liana <AnnaMartinezpisea@lionspest.com.au>
Title : In any case, I am happy that we met
Date  : 2021-06-13 20:57
---
( https://fectment.page.link/Ymry )

What such a gorgeous man is doing here?

*s ho Dent blink scorn league rose ivy superman atkins atkins mugsy freeze thorne katana bane jason edward batarang alfred rumor edward. w ph Maxie vale bartok selina hangman batman young hugo knight freeze batgirl ragman jason batmobile fairchild mister grayson ghul solomon the. ot Elongated czonk diamond bennett batmobile martha hatter snake bruce swamp strange blink creeper abattoir flash sinestro falcone harley bane ragdoll. o* ( https://fectment.page.link/CF1b )

( https://matering.page.link/bAmq )

Will you come to me on the weekend?

*s ho Todd aquaman bullock falcone jester chase croc doom swamp sinestro hangman fairchild nocturna hangman creeper hangman caird aquaman kane barrow. w p Clench chill green canary metallo face robin shrike hatter riddler gleeson justice rumor batarang kane lucius ragman fox grey batmobile. ho Night gleeson oswald cluemaster abattoir ragman gleeson oswald elongated batmobile face quinn abbott clayface moth knight prey knight atkins killer? to* ( https://exteleer.page.link/kjcS )
---
`, content)

	mail, err = Parse[client.MailSourceDoc](getDoc[client.MailSourceDoc](t, "source_mail.html"))
	assert.NoError(t, err)

	content, err = mail.Coloured()
	assert.NoError(t, err)
	assert.Equal(t, `---
Content-Transfer-Encoding : quoted-printable

Content-Type              : text/html; charset=utf-8

Date                      : Sun, 13 Aug 2023 22:45:09 +0000

Dkim-Signature[0]         : v=1; a=rsa-sha256; q=dns/txt; c=relaxed/simple; s=q5
                            bw7xixmmvalostlj63tyl4baejvbto; d=notificacionesatla
                            s.com; t=1691966709; h=Sender:Message-ID:Date:Subjec
                            t:From:To:MIME-Version:Content-Type:Content-Transfer
                            -Encoding; bh=Zl0o2FHSd18g28sSGzovn6Xq/HWn9YPl2DFr+D
                            d+KE4=; b=DTOkYKD3HyTGHUOTGRGL0V2nTOPes9HlNBXHcSms0X
                            Hdr7xL1AXriMYTLwuv1UTM 5iO0ZTFPpMQDfjd7mi/Ca0oNVUAmg
                            aSojcuxWUHu5znCt3e3OSEL8q5u9rN5fI3jFkj ASRgVFTIvJNhH
                            17o44ONqwpIdt2cYd17LMBAfp1f4KK9lPERd0H2jX8SIjc4dHEQx
                            a5 5JDAQN92SlVV6CkhcZYF2mdEhsYuZsPkFVSd6BKlKNPT2Y4tZ
                            iEW5lI+UjTvvbdlRWj i/7ATftL+CYE/mz7soGeeJXV+PNKX4Mgb
                            z8jujp2nV/PrJlZSp7IijF3K/piMTV4udN 6yG/+O1V+Q==

Dkim-Signature[1]         : v=1; a=rsa-sha256; q=dns/txt; c=relaxed/simple; s=22
                            4i4yxa5dv7c2xz3womw6peuasteono; d=amazonses.com; t=1
                            691966709; h=Sender:Message-ID:Date:Subject:From:To:
                            MIME-Version:Content-Type:Content-Transfer-Encoding:
                            Feedback-ID; bh=Zl0o2FHSd18g28sSGzovn6Xq/HWn9YPl2DFr
                            +Dd+KE4=; b=aIwZk+y/naOdqtrYzyFrc8/qkfwgJt6APQ6vP22z
                            qLe5/oLJ23M1KFTbyKCqXlKF t4W1TktUHy2iGXzZB3izHAFHmPA
                            ZmvaplA59iYQsGQI38bZNhf8Dsczpugwm/zy/hTX 7q2ZNub78+g
                            qsXoaoyTSPOcdFhwFrlSfbvxZ14bo=

Feedback-Id               : 1.us-east-1.kRR7d+JzqofruPoUpbLTHFnCtNSHgd8N+6f35f6u
                            eyg=:AmazonSES

From                      : Ola no-reply <aplicativos@notificacionesatlas.com>

Message-Id                : <01000189f1131ee1-caa9ba5a-7352-4f31-a033-df29983a54
                            cc-000000@email.amazonses.com>

Mime-Version              : 1.0

Sender                    : Ola no-reply <aplicativos@notificacionesatlas.com>

Subject                   : =?utf-8?Q?Marcaci=C3=B3n?= de un punto de ronda fuer
                            a de la =?utf-8?Q?posici=C3=B3n?= georreferencia del
                             cliente en INTERCOLOMBIA S.A. E.S.P., zona: RONDA C
                            ASA FEISA.

To                        : test@yopmail.com

X-Ses-Outgoing            : 2023.08.13-54.240.48.111

---
<p>Hola,<br />
<br />
Marcaci=C3=B3n de un punto de ronda fuera d=
e la posici=C3=B3n georreferencia del cliente:<br />
<b>Fecha y hora d=
e la ronda</b>: 2023-08-13 17:15:00<br />
<b>Responsable asignado</b=
>: <br />
<b>Jornada</b>: 24 HORAS > DIURNA<br />
<b>Cliente</b>=
: ISA INTERCOLOMBIA SA ESP<br />
<b>Sede o Punto del Cliente</b>: IN=
TERCOLOMBIA S.A. E.S.P.<br />
<b>Zona interna</b>: RONDA CASA FEISA<=
br />=20
<b>Coordenadas del cliente</b>: 6.1870833;-75.5596067<br =
/>=20
<b>Radio georreferenciado para validar las marcaciones se realice=
n dentro de dicha geocerca (metros)</b>: <br />=20
<b>Distancia de la =
marcaci=C3=B3n respecto a la sede (metros)</b>: 330.06859458703<br />=
=20
</p>."

---
`, content)
}
