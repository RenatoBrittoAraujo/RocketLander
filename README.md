# RocketLander
2D Go simulation about landing a rocket SpaceX style using AI

<img align="center" width="900" height="500" src="https://media.giphy.com/media/JUwr9tx3HGn8q0Ws2t/giphy.gif">

## Links

[Trello Backlog](https://trello.com/b/T70DyPsA/rocketlander)


## Running

```
go install;$GOBIN/rl [args]
```

#### Rocket Lander

Um projeto em Go para simular um foguete pousando estilo SpaceX usando inteligência artificial. O projeto permite que um algoritmo hardcoded, input de usuário e inteligência artificial controle o foguete.

##### Primeira versão do projeto

O foguete se caracteriza por um retângulo com um lado 10 vezes mais comprido que o outro, sendo que um dos lados mais compridos é o topo do foguete e outro a base. Da base do foguete é de onde sai seu impulso (seus motores). O foguete também pode ser controlado com pequenos jatos saindo do seu topo que alteram a direção para onde o foguete vira. Os motores do foguete podem ter seu impulso controlado, variando de 100% até 40% ([Merlin 1D](https://en.wikipedia.org/wiki/SpaceX_Merlin#Merlin_1D)). Inicialmente decola para cima e vira levemente para o lado (aleatório). Após isso, a simulação libera o input para o controlador, seja o computador ou usuário. O resultado final de um pouso é:

1) O foguete entrou em contato com o chão


2) A velocidade (horizontal + vertical) do foguete ao tocar no chão não pode exceder 20 m/s.


3) O vetor que corre do topo do foguete até a base deve fazer ângulo de - 90° +- 5° com o vetor que corre da esquerda até a direita do plano de coordenadas.

4) Após tocar no chão, o input do usuário é cortado e quatro segundos são contados em que se checa se o ângulo citado no item anterior permanece dentro dos limites.

Uma nota é dada para o desempenho de pouso. Ela é positiva se e apenas se todos os itens anteriores forem atendidos. Seu valor é proporcional à proximidade do ângulo de pouso ao ângulo desejado e inversamente proporcional à velocidade de pouso.

Existem três níveis de simulação, variando da mais simples à mais complexa:

1) O foguete tem combustível limitado, pode ligar e desligar seus motores indefinidamente, pode pousar em qualquer lugar.


2) O foguete tem combustível limitado, pode ligar seus motores apenas 2 vezes (sendo uma delas na decolagem e consequente desligamento antes de passar o controle para o computador/usuário), pode pousar em qualquer lugar.


3) O foguete tem combustível limitado, pode ligar seus motores apenas 2 vezes (sendo uma delas na decolagem e consequente desligamento antes de passar o controle para o computador/usuário), deve pousar dentro de área específica delimitada no eixo de coordenadas X (ou seja, nenhum ponto do retângulo que o cobre pode estar fora dos limites mínimos e máximos dessa área). Neste caso em específico, além da pontuação ser dada pelo ângulo e velocidade de toque no chão, ela é inversamente proporcional à distância do centro do foguete à região de pouso.

O projeto deve ser desenvolvido em Golang, usando a biblioteca gráfica Ebiten. O desenvolvimento vai ser decomposto em etapas:

1) A simulação deve ser criada, rodando a exatamente 30 frames por segundo seguindo as especificações anteriores. Á cada novo frame, um objeto de input (podendo ser esse: algoritmo hardcoded, algoritmo de inteligência artificial ou input direto de usuário) receberá os dados relevantes do foguete

 - sua posição
 - ângulo em relação ao chão
 - vetor de velocidade
 - combustível restante
 - status do motor (ligado ou não)
 - porcentagem de impulso do motor.
 - status dos motores laterais (esquerda, direita ou nenhum)

  E retornará nenhuma ou várias mudanças sobre os seguintes dados:

 - status do motor
 - porcentagem de impulso do motor
 - status do motores laterais

  Após essas mudanças, um novo frame será recalculado e o processo se repete até que o foguete toque no chão.

2) Uma rede neural deve ser criada sem uso de bibliotecas, com uma camada de input suportando os inputs acima e uma camada de output também suportando os dados acima. Sua rede deve suportar um número variável de camadas do meio e, para cada camada, um número variável de neurônios.

3) A adição de algoritmo hardcoded é opcional, mas deve ser fácil de ligar ao programa.

Notas adicionais:

- Como a rede neural deve treinar, para acelerar esse processo, deve ser possível aumentar a quantidade de frames que ocorrem por segundo pela própria linha de comando ao iniciar o programa. As mudanças entre frames devem ser constantes, ou seja, independente de quantos frames por segundo o jogo está rodando, o impacto resultante sobre a física dos objetos no jogo deve ser a mesma.

- Deve ser possível rodar a simulação sem desenhar na tela o que está ocorrendo. Uma opção de linha de comando deve ser adicionada para rodar a simulação e no terminal mostrar seu progresso e resultados.
