# 📱 Estratégia: Responsive SOTA & Mobile First (CSS Puro)

Esta diretriz define como o Arandu evolui para ser onipresente na rotina do terapeuta, utilizando **CSS Nativo** e `templ`, mantendo a integridade da "Tecnologia Silenciosa".

## 1. O Racional: Porquê Mobile First?

1. **Hierarquia de Informação:** O espaço reduzido obriga a decidir o que é vital (ex: relato da sessão) versus o que é secundário.

2. **Ergonomia (A Lei do Polegar):** Em dispositivos móveis, a zona inferior é a mais nobre. Acções críticas devem estar na Bottom Nav.

3. **Adaptação Fluida:** O layout expande-se para revelar camadas contextuais no desktop através de Media Queries.

## 2. Matriz de Navegação Dual

| Elemento | Comportamento Mobile (< 768px) | Comportamento Desktop (>= 768px) | 
| ----- | ----- | ----- | 
| **Sidebar** | **Drawer:** Oculto por padrão. Acionado por ícone. | **Persistente:** Fixa à esquerda (`280px`). | 
| **Ações Rápidas** | **Bottom Nav:** Fixa no rodapé para Navegação Principal. | **Sidebar/Header:** Integrado no menu lateral. | 
| **Timeline** | **Full Width:** Empilhamento vertical. | **Lateral:** Flutua à direita do editor. | 

## 3. Diretrizes Técnicas (CSS Nativo)

### A. Navegação Mobile (Bottom Bar)

Deve ser implementada no `base.css` com visibilidade alternada:

```css
.bottom-nav {
    position: fixed;
    bottom: 0;
    display: flex;
}
@media (min-width: 768px) {
    .bottom-nav { display: none; }
}

B. O Drawer (Sidebar Mobile)
O componente aside deve usar transições de transform:

.sidebar-drawer: fixed, z-index: 50, transform: translateX(-100%).

.sidebar-open: transform: translateX(0).

🛡️ Checklist de Validação Responsiva
[ ] O componente é legível a 320px de largura?

[ ] Os botões principais estão na "zona do polegar"?

[ ] A sidebar recolhível (Drawer) fecha ao clicar fora (Alpine.js @click.away)?

[ ] O conteúdo clínico mantém o tamanho text-xl (aprox. 18-20px) para evitar zoom.