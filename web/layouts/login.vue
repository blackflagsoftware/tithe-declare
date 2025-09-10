<template>
	<div class="outer">
		<div class="inner">
			<slot />
		</div>
	</div>
</template>

<style>
.outer {
	top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;	
}

.inner {
	display: flex;
  flex-direction: column;
  background-color: var(--color-secondary-bg);
  color: var(--color-secordary-fg);
  padding: 30px;
  border-radius: 10px;
  width: 400px;
}
</style>